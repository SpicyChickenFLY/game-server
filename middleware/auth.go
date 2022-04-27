package middleware

import (
	"net/http"
	"strings"

	"github.com/SpicyChickenFLY/game-server/utils"
	"github.com/gin-gonic/gin"
)

// AuthJWT 验证 JWT 的中间件
func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		headerList := strings.Split(header, " ")
		if len(headerList) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "无法解析 Authorization 字段"})
			c.Abort()
			return
		}
		t := headerList[0]
		content := headerList[1]
		if t != "Bearer" {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "认证类型错误, 当前只支持 Bearer"})
			c.Abort()
			return
		}
		if _, err := utils.Verify([]byte(content)); err != nil {
			c.Abort()
			return
		}

		c.Next()
	}
}
