package user

type User struct {
	UserId   uint64 `json:"user_id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

type RoleEnum string

const (
	USER  RoleEnum = "USER"
	ADMIN RoleEnum = "ADMIN"
	TEST  RoleEnum = "TEST"
)

type Role struct {
	RoleId uint64   `json:"role_id"`
	Role   RoleEnum `json:"role"`
	UserId uint64   `json:"user_id"`
}

type UserDto struct {
	Fullname string   `json:"fullname"`
	Email    string   `json:"email"`
	Role     RoleEnum `json:"role"`
}
