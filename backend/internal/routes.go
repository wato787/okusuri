package internal

import (
	"okusuri-backend/internal/handler"
	"okusuri-backend/internal/middleware"
	"okusuri-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	// リポジトリの初期化（Cognitoベース）
	medicationRepo := repository.NewMedicationRepository()
	notificationRepo := repository.NewNotificationRepository()

	// ハンドラーの初期化
	medicationHandler := handler.NewMedicationHandler(medicationRepo)
	notificationHandler := handler.NewNotificationHandler(notificationRepo)

	// Ginのルーターを作成
	router := gin.Default()

	// グローバルミドルウェアの設定
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// Cognito認証必須エンドポイント
		api.GET("/medication-status", middleware.CognitoAuth(), medicationHandler.GetMedicationStatus)

		medicationLog := api.Group("/medication-log")
		medicationLog.Use(middleware.CognitoAuth())
		{
			medicationLog.POST("", medicationHandler.RegisterLog)
			medicationLog.GET("", medicationHandler.GetLogs)
			medicationLog.GET("/:id", medicationHandler.GetLogByID)
			medicationLog.PATCH("/:id", medicationHandler.UpdateLog)
		}

		// 通知設定エンドポイント
		notificationSetting := api.Group("/notification/setting")
		notificationSetting.Use(middleware.CognitoAuth())
		{
			notificationSetting.GET("", notificationHandler.GetSetting)
			notificationSetting.POST("", notificationHandler.RegisterSetting)
		}
	}

	return router
}
