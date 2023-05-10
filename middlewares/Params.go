package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"unicode"
)

func FilterParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求参数
		params := c.Request.URL.Query()
		// 遍历参数，检查是否只包含数字和字母
		for key, values := range params {
			for _, value := range values {
				if !isAlphanumeric(value) {
					// 如果不是数字和字母，返回错误
					c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid parameter %s: %s", key, value)})
					c.Abort() // 终止后续处理
					return
				}
			}
		}
		// 如果参数都合法，继续后续处理
		c.Next()
	}
}

func isAlphanumeric(input string) bool {
	for _, r := range input {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
