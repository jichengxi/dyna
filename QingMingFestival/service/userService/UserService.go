package userService

import (
	"QingMingFestival/ormOperate/userOperate"
	"QingMingFestival/types/userTypes"
	"strings"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) GetUserInfo() (*userTypes.RequestUserInfo, error) {
	user, err := userOperate.NewUserOperate().QueryUser()
	if err != nil {
		return nil, err
	}
	role := strings.Split(user.Role, ",")
	return &userTypes.RequestUserInfo{
		Id:     user.Id,
		Name:   user.Name,
		Avatar: user.Avatar,
		Email:  user.Email,
		Role:   role,
	}, nil
}
