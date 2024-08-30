package user

import (
	"context"
	"database/sql"
	"github.com/codescratchers/golang-webserver/database"
	"log"
	"reflect"
	"testing"
)

var dbInstance *sql.DB

func TestMain(m *testing.M) {
	db, err := database.ConnectToMySQL("hushtabs", "hushtabs", "127.0.0.1:3306", "hushtabs_db")

	if err != nil {
		log.Fatal(err)
	}
	dbInstance = db

	// close db connection
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Print("db connection did not close after tests")
			return
		}
	}(db)

	// run tests
	m.Run()
}

func SetupTest(t *testing.T) *sql.Tx {
	tx, err := dbInstance.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
		return nil
	}

	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	})

	return tx
}

func TestSave_And_UserByEmail(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("TestGivenEmailThatExist_ShouldReturnSaidUser", func(t *testing.T) {
		t.Parallel()
		tx := SetupTest(t)

		user := User{
			Fullname: "Test User",
			Email:    "demo@email.com",
		}

		if err := saveTx(ctx, tx, &user); err != nil {
			t.Errorf("%s", err)
		}

		// method to test
		find, err := userByEmailTx(ctx, tx, user.Email)

		// assert
		if err != nil {
			t.Errorf("%s", err)
		} else if !reflect.DeepEqual(user, find) {
			t.Errorf("assertion failed user save not equal to user found")
		}
	})

	t.Run("TestGivenEmail_ShouldReturnNoUser", func(t *testing.T) {
		t.Parallel()
		tx := SetupTest(t)

		// method to test
		_, err := userByEmailTx(ctx, tx, "demo@email.com")

		// assert
		if err == nil {
			t.Errorf("user should not be found %s", err)
		}
	})
}
