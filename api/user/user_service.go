package user

import "context"

type IUserService interface {
	CreateUser(ctx context.Context, fullname string) (User, error)
}

type userService struct {
	UserRepository IUserRepository
}

// NewUserService constructor
func NewUserService(r IUserRepository) IUserService {
	return userService{UserRepository: r}
}

// CreateUser saves a User obj to the database if the email does not exist
func (s userService) CreateUser(ctx context.Context, fullname string) (User, error) {
	return s.UserRepository.Save(ctx, User{Fullname: fullname})
}
