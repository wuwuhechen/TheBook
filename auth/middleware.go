package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未提供认证Token",
			})
			c.Abort()
			return
		}

		if len(token) < 7 || token[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token格式错误",
			})
			c.Abort()
			return
		}

		if len(token) >= 7 {
			token = token[7:]
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token无效或已过期",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.UserName)
		c.Set("nickname", claims.Nickname)
	}
}
