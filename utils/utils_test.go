package utils

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestDbTransaction(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("GivenFunctionReturnsError_ShouldRollbackTransaction", func(t *testing.T) {
		t.Parallel()

		// given
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("db not initialized")
		}
		defer db.Close()

		mock.ExpectBegin()
		tx, err := db.BeginTx(ctx, nil)

		if err != nil {
			t.Errorf("tx exception")
		}

		mock.ExpectRollback()

		function := func(ctx context.Context) error {
			return fmt.Errorf("simulate db execption")
		}

		// method to test and assert
		if err := DbTransaction(ctx, tx, function); err == nil {
			t.Errorf("error was not returned so rollback was not initiated")
		}
	})

	t.Run("GivenFunctionThatDoesn'tReturnError_ShouldCommitTransaction", func(t *testing.T) {
		t.Parallel()

		// given
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("db not initialized")
		}
		defer db.Close()

		mock.ExpectBegin()
		tx, err := db.BeginTx(ctx, nil)

		if err != nil {
			t.Errorf("tx exception")
		}

		mock.ExpectCommit()

		function := func(ctx context.Context) error {
			return nil
		}

		// method to test and assert
		if err := DbTransaction(ctx, tx, function); err != nil {
			t.Errorf("rollback was not initiated")
		}
	})
}
