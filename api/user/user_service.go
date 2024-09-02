package user

import (
	"context"
	"database/sql"
	"github.com/codescratchers/golang-webserver/database"
)

type IUserService interface {
	UserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, dto UserDto) error
}

type userService struct {
	db             *sql.DB
	UserRepository IUserRepository
	RoleRepository IRoleRepository
}

// NewUserService constructor
func NewUserService(db *sql.DB, u IUserRepository, r IRoleRepository) IUserService {
	return userService{db: db, UserRepository: u, RoleRepository: r}
}

func (s userService) UserByEmail(ctx context.Context, email string) (User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return User{}, err
	}
	var user User
	err = database.DbTransaction(ctx, tx, func(ctx context.Context) error {
		user, err = s.UserRepository.UserByEmail(tx, email)
		return err
	})
	return user, err
}

func (s userService) CreateUser(ctx context.Context, dto UserDto) error {
	// if email found duplicate user
	if _, err := s.UserByEmail(ctx, dto.Email); err == nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	return database.DbTransaction(ctx, tx, func(ctx context.Context) error {
		user := User{Fullname: dto.Fullname, Email: dto.Email}
		if err := s.UserRepository.Save(tx, &user); err != nil {
			return err
		}
		return s.RoleRepository.Save(tx, Role{Role: dto.Role, UserId: user.UserId})
	})
}
