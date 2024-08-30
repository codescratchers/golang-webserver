package database

import (
	"context"
	"database/sql"
	"fmt"
)

func DbTransaction(ctx context.Context, tx *sql.Tx, f func(context.Context) error) error {
	if err := f(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("db rollback %w", err)
		}
		return fmt.Errorf("f %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db commit %w", err)
	}

	return nil
}
