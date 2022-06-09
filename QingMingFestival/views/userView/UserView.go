package userView

import (
	"QingMingFestival/service/userService"
	"QingMingFestival/tools"
	"QingMingFestival/tools/middleware"
	"QingMingFestival/types/userTypes"
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserView struct {
}

func (uv *UserView) Router(api *gin.RouterGroup) {
	api.GET("/user", uv.getUserInfo)
	api.POST("/auth/login", uv.login)
	api.POST("/auth/refreshToken", uv.refreshToken)
}

func (uv *UserView) getUserInfo(ctx *gin.Context) {
	userInfoRequest, err := userService.NewUserService().GetUserInfo()
	if err != nil {
		middleware.Logf.Error(err)
		tools.Failed(ctx, "获取用户信息失败.")
		return
	}
	middleware.Logf.Info("获取用户信息成功.")
	tools.Success(ctx, *userInfoRequest)
}

type testToken struct {
	Token string `json:"token"`
}

func (uv *UserView) login(ctx *gin.Context) {
	var userAuth userTypes.UserAuth
	err := tools.JsonDecode(ctx.Request.Body, &userAuth)
	fmt.Println(userAuth)
	if err != nil {
		middleware.Logf.Error(err)
		tools.Failed(ctx, "用户信息参数有误.")
		return
	}
	fmt.Println("userAuth: ", userAuth)
	// todo: 测试 随便返回
	tools.Success(ctx, testToken{Token: "admin"})
}

func (uv *UserView) refreshToken(ctx *gin.Context) {
	tools.Success(ctx, testToken{Token: "admin"})
}
