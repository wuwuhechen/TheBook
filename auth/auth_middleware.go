package auth

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 验证请求中的 JWT，并将用户信息写入 Gin Context。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			cookieToken, err := c.Cookie("access_token")
			if err != nil || cookieToken == "" {
				redirectToLogin(c, "未提供认证 Token")
				return
			}
			token = cookieToken
		}

		if len(token) >= 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := ParseToken(token)
		if err != nil {
			redirectToLogin(c, "Token 无效或已过期")
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.UserName)
		c.Set("nickname", claims.Nickname)
		c.Next()
	}
}

func redirectToLogin(c *gin.Context, message string) {
	c.Redirect(http.StatusFound, "/auth/login?message="+url.QueryEscape(message))
	c.Abort()
}
