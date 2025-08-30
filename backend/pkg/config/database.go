package config

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/guregu/dynamo/v2"
	"github.com/rs/zerolog/log"
)

// DB はアプリケーション全体で使用するDynamoDB接続です
var DB *dynamo.DB

// SetupDB はDynamoDB接続を初期化します
func SetupDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info().Msg("DynamoDB接続を初期化しています")

	// AWS設定を読み込み
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("AWS設定の読み込みに失敗しました")
	}

	// 環境変数の確認
	region := cfg.Region
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")

	log.Info().
		Str("region", region).
		Str("endpoint", endpoint).
		Msg("AWS設定を読み込みました")

	// DynamoDB接続を初期化
	db := dynamo.New(cfg)

	log.Info().Msg("DynamoDB接続に成功しました")
	DB = db
}

// GetDB はDynamoDB接続を返します
func GetDB() *dynamo.DB {
	if DB == nil {
		log.Fatal().Msg("DynamoDB接続が初期化されていません。SetupDB()を先に呼び出してください")
	}
	return DB
}
