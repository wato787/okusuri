package middleware

import (
	"okusuri-backend/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CognitoAuth はAPI GatewayのCognito認証から渡されるユーザーIDを検証するミドルウェア
func CognitoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// API GatewayのCognito認証から渡されるユーザーID
		userID := c.Request.Header.Get("X-Cognito-User-Id")
		if userID == "" {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("user_agent", c.Request.UserAgent()).
				Msg("Cognito user ID not found in request header")

			errors.HandleUnauthorized(c, "認証が必要です", nil, "Cognito user ID not found")
			c.Abort()
			return
		}

		// ユーザーIDの基本検証
		if len(userID) < 10 || len(userID) > 128 {
			log.Warn().
				Str("user_id", userID).
				Str("path", c.Request.URL.Path).
				Msg("Invalid Cognito user ID format")

			errors.HandleUnauthorized(c, "無効なユーザーIDです", nil, "Invalid user ID format")
			c.Abort()
			return
		}

		log.Debug().
			Str("user_id", userID).
			Str("path", c.Request.URL.Path).
			Msg("Cognito authentication successful")

		// Cognitoユーザー情報をコンテキストに保存
		c.Set("cognitoUserID", userID)

		c.Next()
	}
}
