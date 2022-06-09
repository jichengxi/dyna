package main

import (
	"QingMingFestival/ormOperate/vmOperate"
	"QingMingFestival/tools"
	"QingMingFestival/tools/middleware"
	"QingMingFestival/views"
	"QingMingFestival/views/userView"
	"QingMingFestival/views/vmView"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := tools.ParseConfig("./config/app.json")
	if err != nil {
		panic(err)
	}
	tools.OrmEngine(cfg)

	engine := gin.Default()
	// cors中间件
	engine.Use(middleware.Cors(), middleware.LoggerToFile())

	// 全局路由组
	api := engine.Group("/api")
	registerRouter(api)

	engine.Run(cfg.AppHost + ":" + cfg.AppPort)
	defer func() {
		oldClient := vmOperate.GetVmClient()
		if oldClient != nil {
			oldClient.LogoutVc()
		}
	}()
}

// 路由设置
func registerRouter(apiGroup *gin.RouterGroup) {
	new(views.HelloView).Router(apiGroup)
	new(vmView.VMView).Router(apiGroup)
	new(userView.UserView).Router(apiGroup)
}
