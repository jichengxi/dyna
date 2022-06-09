package userOperate

import (
	"QingMingFestival/model/userModel"
	"QingMingFestival/tools"
	"gorm.io/gorm"
)

type UserOperate struct {
	*gorm.DB
}

func NewUserOperate() *UserOperate {
	return &UserOperate{tools.DbEngine}
}

func (uo *UserOperate) QueryUser() (*userModel.UserInfo, error) {
	var userInfos userModel.UserInfo
	res := uo.Where("role = ?", "admin").First(&userInfos)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userInfos, nil
}
