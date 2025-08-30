package main

import (
	routes "okusuri-backend/internal"
	"okusuri-backend/pkg/config"
	"okusuri-backend/pkg/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	// ログ初期化
	logger.InitLogger()

	log.Info().Msg("アプリケーションを開始します")

	// DynamoDB接続
	config.SetupDB()
	log.Info().Msg("DynamoDB接続が完了しました")

	// Ginのルーターを作成
	router := routes.SetupRoutes()
	log.Info().Msg("ルーター設定が完了しました")

	// ポート設定（Lambda Web Adapter対応）
	port := config.GetPort()

	log.Info().Str("port", port).Msg("サーバーを起動します")

	// サーバーを起動
	if err := router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("サーバーの起動に失敗しました")
	}
}
