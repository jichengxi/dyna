package userTypes

type RequestUserInfo struct {
	Id     int64    `json:"id"`
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
	Email  string   `json:"email"`
	Role   []string `json:"role"`
}
