package views

import "github.com/gin-gonic/gin"

type HelloView struct {

}

func (hello *HelloView) Router(api *gin.RouterGroup)  {
	api.GET("/hello", hello.Hello)
}

// Hello 解析 /hello 请求
func (hello *HelloView) Hello(c *gin.Context)  {
	c.JSON(200, map[string]interface{}{
		"message": "hello 清明节",
	})
}

