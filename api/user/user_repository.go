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
	UserByEmail(ctx context.Context, email string) (User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return userRepository{db: db}
}

func (r userRepository) UserByEmail(ctx context.Context, email string) (User, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return User{}, fmt.Errorf("tx begin %w", err)
	}

	var user User

	err = database.DbTransaction(ctx, tx, func(context.Context) error {
		if u, err := userByEmailTx(ctx, tx, email); err != nil {
			return err
		} else {
			user = u
			return nil
		}
	})

	if err != nil {
		return User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func userByEmailTx(ctx context.Context, tx *sql.Tx, email string) (User, error) {
	var user User
	row := tx.QueryRowContext(ctx, "SELECT * FROM user WHERE email = ?", email)

	// scan the result into the User struct
	if err := row.Scan(
		&user.UserId,
		&user.Fullname,
		&user.Email,
	); err != nil {
		return User{}, err
	}

	return user, nil
}

func (r userRepository) Save(ctx context.Context, user User) (User, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return User{}, fmt.Errorf("save user tx begin %w", err)
	}

	var u = user
	err = database.DbTransaction(ctx, tx, func(context.Context) error {
		return saveTx(ctx, tx, &u)
	})

	return u, err
}

func saveTx(ctx context.Context, tx *sql.Tx, user *User) error {
	result, err := tx.ExecContext(
		ctx,
		"INSERT INTO user (fullname, email) VALUES (?, ?)",
		user.Fullname,
		user.Email,
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
