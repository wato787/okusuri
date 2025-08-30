package middleware

import (
	"github.com/gin-gonic/gin"
)

// CognitoAuth はAPI GatewayのCognito認証から渡されるユーザーIDを検証するミドルウェア
func CognitoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// API GatewayのCognito認証から渡されるユーザーID
		userID := c.Request.Header.Get("X-Cognito-User-Id")
		if userID == "" {
			c.JSON(401, gin.H{"error": "Cognito user ID not found in request header"})
			c.Abort()
			return
		}

		// Cognitoユーザー情報をコンテキストに保存
		c.Set("cognitoUserID", userID)

		c.Next()
	}
}
