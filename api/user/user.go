package user

type User struct {
	UserId   uint64 `json:"user_id"`
	Fullname string `json:"fullname"`
}

type UserDto struct {
	Fullname string `json:"fullname"`
}
