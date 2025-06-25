package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 日志中间件
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 处理请求
		c.Next()
		// 跳过健康检查的日志
		if c.Request.URL.Path == "/health" {
			return
		}
		// 获取请求信息
		method := c.Request.Method
		url := c.Request.URL.String()
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		// 获取请求体（限制大小）
		var body string
		if c.Request.Body != nil {
			bodyBytes, _ := c.GetRawData()
			if len(bodyBytes) > 0 {
				if len(bodyBytes) > 1024*1024 {
					body = "body is too large"
				} else {
					body = string(bodyBytes)
				}
			}
		}
		// 获取响应体（限制大小）
		var result string
		if c.Writer.Size() > 0 {
			if c.Writer.Size() > 1024*1024 {
				result = "result is too large"
			} else {
				// 注意：gin 不直接提供响应体内容，这里只是记录大小
				result = fmt.Sprintf("response size: %d", c.Writer.Size())
			}
		} // 记录日志
		zap.L().Info(c.Request.URL.Path,
			zap.String("method", method),
			zap.String("url", url),
			zap.String("client", clientIP),
			zap.String("body", body),
			zap.Int("code", statusCode),
			zap.String("result", result),
			zap.Duration("latency", time.Since(start)))
	}
}
