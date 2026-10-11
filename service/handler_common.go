package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlerGetHomePage 渲染应用首页。
func (s *Server) HandlerGetHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home_page.html", nil)
}

// currentUserID 从认证中间件写入的 Context 中读取用户主键。
func currentUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := value.(uint)
	return id, ok
}

// choiceText 将选项索引格式化为带字母前缀的文本。
func choiceText(choices []string, index int) string {
	if index < 0 || index >= len(choices) {
		return ""
	}
	return string(rune('A'+index)) + ". " + choices[index]
}
