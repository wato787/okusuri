package config

import (
	"os"
)

// Environment はアプリケーションの環境変数を管理します
type Environment struct {
	// サーバー設定
	Port string

	// AWS設定
	AWSRegion string

	// DynamoDB設定
	DynamoDBTableName string

	// ログ設定
	LogLevel string
}

// Load は環境変数から設定を読み込みます
func Load() *Environment {
	return &Environment{
		// サーバー設定
		Port: getEnv("PORT", "8080"),

		// AWS設定
		AWSRegion: getEnv("AWS_REGION", "ap-northeast-1"),

		// DynamoDB設定
		DynamoDBTableName: getEnv("DYNAMODB_TABLE_NAME", "okusuri-production-table"),

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

// GetPort はサーバーポートを取得します
func GetPort() string {
	return Load().Port
}

// GetLogLevel はログレベルを取得します
func GetLogLevel() string {
	return Load().LogLevel
}
