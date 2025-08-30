package handler

import (
	"net/http"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"okusuri-backend/pkg/helper"
	"time"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationHandler(notificationRepo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
	}
}

// GetSetting はユーザーの通知設定を取得するハンドラー
func (h *NotificationHandler) GetSetting(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// クエリパラメータからプラットフォームを取得
	platform := c.DefaultQuery("platform", "web")

	// 通知設定を取得
	setting, err := h.notificationRepo.GetSetting(userID, platform)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification setting not found"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// RegisterSetting はユーザーの通知設定を登録/更新するハンドラー
func (h *NotificationHandler) RegisterSetting(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// リクエストボディを構造体にバインド
	var req dto.NotificationSettingRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// NotificationSetting構造体を作成
	setting := model.NotificationSetting{
		Platform:     req.Platform,
		IsEnabled:    req.IsEnabled,
		Subscription: req.Subscription,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// リポジトリを使って保存
	err = h.notificationRepo.RegisterSetting(userID, setting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register notification setting"})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "notification setting registered successfully",
	})
}
