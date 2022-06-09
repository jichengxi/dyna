package tools

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SUCCESS = 0
	FAILED  = 1
)

// Success 普通成功返回
func Success(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    SUCCESS,
		"message": "成功",
		"data":    v,
	})
}

func Failed(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    FAILED,
		"message": v,
		"data":    "",
	})
}
