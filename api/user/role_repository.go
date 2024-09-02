package user

import (
	"context"
	"database/sql"
	"github.com/codescratchers/golang-webserver/database"
)

type IRoleRepository interface {
	Save(ctx context.Context, role Role) error
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) IRoleRepository {
	return roleRepository{db: db}
}

func (r roleRepository) Save(ctx context.Context, role Role) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	return database.DbTransaction(ctx, tx, func(context.Context) error {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO role (role, user_id) VALUES (?, ?)",
			role.Role,
			role.UserId,
		)
		return err
	})
}
