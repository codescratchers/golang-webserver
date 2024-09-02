package user

import "context"

type IUserService interface {
	UserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, dto UserDto) (User, error)
}

type userService struct {
	UserRepository IUserRepository
	RoleRepository IRoleRepository
}

// NewUserService constructor
func NewUserService(u IUserRepository, r IRoleRepository) IUserService {
	return userService{UserRepository: u, RoleRepository: r}
}

func (s userService) UserByEmail(ctx context.Context, email string) (User, error) {
	return s.UserRepository.UserByEmail(ctx, email)
}

func (s userService) CreateUser(ctx context.Context, dto UserDto) (User, error) {
	user, err := s.UserRepository.Save(ctx, User{Fullname: dto.Fullname, Email: dto.Email})
	if err != nil {
		return User{}, err
	}
	return user, s.RoleRepository.Save(ctx, Role{Role: dto.Role, UserId: user.UserId})
}
