package user

type User struct {
	UserId   uint64 `json:"user_id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

type UserDto struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}
