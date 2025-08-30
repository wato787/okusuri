package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/guregu/dynamo/v2"
)

// DB はアプリケーション全体で使用するDynamoDB接続です
var DB *dynamo.DB

// SetupDB はDynamoDB接続を初期化します
func SetupDB() {
	// AWS設定を読み込み
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("AWS設定の読み込みに失敗しました: ", err)
	}

	// DynamoDB接続を初期化
	db := dynamo.New(cfg)
	log.Println("DynamoDB接続に成功しました")
	DB = db
}

// GetDB はDynamoDB接続を返します
func GetDB() *dynamo.DB {
	return DB
}
