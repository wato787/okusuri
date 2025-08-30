package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext はコンテキストからCognitoユーザーIDを取得
func GetUserIDFromContext(c *gin.Context) (string, error) {
	// Cognito用のユーザーID取得を優先
	if userID, exists := c.Get("cognitoUserID"); exists {
		if id, ok := userID.(string); ok && id != "" {
			return id, nil
		}
	}

	return "", errors.New("CognitoユーザーIDがコンテキストに存在しません")
}
