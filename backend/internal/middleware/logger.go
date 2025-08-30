package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger は構造化ログを使用するGinミドルウェア
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// ヘルスチェックはログをスキップ
		if path == "/api/health" {
			c.Next()
			return
		}

		// リクエスト処理の前
		c.Next()

		// リクエスト処理の後
		latency := time.Since(start)
		status := c.Writer.Status()

		logger := log.Info()
		if status >= 400 {
			logger = log.Error()
		} else if status >= 300 {
			logger = log.Warn()
		}

		logger.
			Str("client_ip", c.ClientIP()).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status_code", status).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Int("body_size", c.Writer.Size()).
			Msg("HTTP request completed")
	}
}
