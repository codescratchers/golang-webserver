package user

import "context"

type IUserService interface {
	UserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, dto UserDto) (User, error)
}

type userService struct {
	UserRepository IUserRepository
}

// NewUserService constructor
func NewUserService(r IUserRepository) IUserService {
	return userService{UserRepository: r}
}

func (s userService) UserByEmail(ctx context.Context, email string) (User, error) {
	return s.UserRepository.UserByEmail(ctx, email)
}

func (s userService) CreateUser(ctx context.Context, dto UserDto) (User, error) {
	return s.UserRepository.Save(ctx, User{Fullname: dto.Fullname, Email: dto.Email})
}
