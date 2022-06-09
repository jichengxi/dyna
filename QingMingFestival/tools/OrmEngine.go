package tools

import (
	"QingMingFestival/model/userModel"
	"QingMingFestival/model/vmModel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbEngine *gorm.DB

func OrmEngine(cfg *Config) {
	// 用户名:密码@tcp(ip:port)/数据库?charset=utf8mb4&parseTime=True&loc=Local
	database := cfg.Database
	conn := database.User + ":" + database.Password + "@tcp(" + database.Host + ":" + database.Port + ")/" + database.DbName + "?charset=" + database.Charset + "&parseTime=True&loc=Local"
	//dsn := "root:root123@tcp(127.0.0.1:3306)/test_gorm?charset=utf8mb4&parseTime=True&loc=Local"
	dbEngine, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 迁移 schema
	if database.Debug {
		dbEngine = dbEngine.Debug()
	}
	dbEngine.AutoMigrate(new(vmModel.VcAccount), new(vmModel.CloneVmJob), new(userModel.UserInfo), new(userModel.Role))
	DbEngine = dbEngine
}
