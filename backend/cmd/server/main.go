package main

import (
	"log"

	routes "okusuri-backend/internal"
	"okusuri-backend/pkg/config"
)

func main() {
	// DynamoDB接続
	config.SetupDB()

	// Ginのルーターを作成
	router := routes.SetupRoutes()

	// サーバーを起動
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
