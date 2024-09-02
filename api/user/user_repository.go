package user

import (
	"database/sql"
	"fmt"
	"log"
)

type IUserRepository interface {
	Save(tx *sql.Tx, user *User) error
	UserByEmail(tx *sql.Tx, email string) (User, error)
}

type userRepository struct{}

func NewUserRepository() IUserRepository {
	return userRepository{}
}

func (userRepository) UserByEmail(tx *sql.Tx, email string) (User, error) {
	var u User
	row := tx.QueryRow("SELECT * FROM user WHERE email = ?", email)

	// scan the result into the User struct
	if err := row.Scan(
		&u.UserId,
		&u.Fullname,
		&u.Email,
	); err != nil {
		return User{}, err
	}

	return u, nil
}

func (userRepository) Save(tx *sql.Tx, u *User) error {
	result, err := tx.Exec("INSERT INTO user (fullname, email) VALUES (?, ?)", u.Fullname, u.Email)
	if err != nil {
		log.Printf("%v", err)
		return fmt.Errorf("error saving u")
	}
	if id, err := result.LastInsertId(); err != nil {
		log.Printf("error saving u with respect to u id : %v", err)
		return fmt.Errorf("error retrieiving u id after saving")
	} else {
		u.UserId = uint64(id)
		return nil
	}
}
