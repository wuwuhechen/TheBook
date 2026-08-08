package logger

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger(log *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		fileds := []zap.Field{
			zap.String("source", "backend"),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("duration", time.Since(start)),
		}

		if len(c.Errors) > 0 {
			fileds = append(fileds, zap.String("errors", c.Errors.String()))
		}

		switch {
		case c.Writer.Status() >= http.StatusInternalServerError:
			log.Error("http request failed", fileds...)
		case c.Writer.Status() >= http.StatusBadRequest:
			log.Warn("http request rejected", fileds...)
		default:
			log.Info("http request processed", fileds...)
		}
	}
}
