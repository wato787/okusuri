package config

import (
	"os"
)

// Environment はnotification層の環境変数を管理します
type Environment struct {
	// AWS設定
	AWSRegion string

	// DynamoDB設定
	DynamoDBTableName string

	// Cognito設定
	CognitoUserPoolID string

	// Push通知設定
	VAPIDPublicKey  string
	VAPIDPrivateKey string

	// ログ設定
	LogLevel string
}

// Load は環境変数から設定を読み込みます
func Load() *Environment {
	return &Environment{
		// AWS設定
		AWSRegion: getEnv("AWS_REGION", "ap-northeast-1"),

		// DynamoDB設定
		DynamoDBTableName: getEnv("DYNAMODB_TABLE_NAME", "okusuri-production-table"),

		// Cognito設定
		CognitoUserPoolID: getEnv("COGNITO_USER_POOL_ID", ""),

		// Push通知設定
		VAPIDPublicKey:  getEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey: getEnv("VAPID_PRIVATE_KEY", ""),

		// ログ設定
		LogLevel: getEnv("LOG_LEVEL", "INFO"),
	}
}

// getEnv は環境変数を取得し、デフォルト値を設定します
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDynamoDBTableName はDynamoDBテーブル名を取得します
func GetDynamoDBTableName() string {
	return Load().DynamoDBTableName
}

// GetCognitoUserPoolID はCognito User Pool IDを取得します
func GetCognitoUserPoolID() string {
	return Load().CognitoUserPoolID
}

// GetVAPIDPublicKey はVAPID公開鍵を取得します
func GetVAPIDPublicKey() string {
	return Load().VAPIDPublicKey
}

// GetVAPIDPrivateKey はVAPID秘密鍵を取得します
func GetVAPIDPrivateKey() string {
	return Load().VAPIDPrivateKey
}
