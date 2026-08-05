package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			cookieToken, err := c.Cookie("access_token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "未提供认证Token",
				})
				c.Abort()
				return
			}
			token = cookieToken
		}

		if len(token) >= 7 && token[:7] == "Bearer " {
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
