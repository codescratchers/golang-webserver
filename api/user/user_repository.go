package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/codescratchers/golang-webserver/database"
	"log"
)

type IUserRepository interface {
	Save(ctx context.Context, user User) (User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return userRepository{db: db}
}

func (r userRepository) Save(ctx context.Context, user User) (User, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return User{}, fmt.Errorf("save user tx begin %w", err)
	}

	i := 0
	var u = user
	err = database.DbTransaction(ctx, tx, func(context.Context) error {
		return saveTx(ctx, tx, &u)
	})

	log.Printf("i %v", i)

	return u, err
}

func saveTx(ctx context.Context, tx *sql.Tx, user *User) error {
	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO user (fullname) VALUES (?)",
		user.Fullname,
	)

	if err != nil {
		log.Printf("%v", err)
		return fmt.Errorf("error saving user")
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Printf("error saving user with respect to user id : %v", err)
		return fmt.Errorf("error retrieiving user id after saving")
	}

	user.UserId = uint64(id)

	return nil
}
