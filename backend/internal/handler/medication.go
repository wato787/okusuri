package handler

import (
	"net/http"
	"okusuri-backend/internal/dto"
	"okusuri-backend/internal/model"
	"okusuri-backend/internal/repository"
	"okusuri-backend/internal/service"
	"okusuri-backend/pkg/errors"
	"okusuri-backend/pkg/helper"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type MedicationHandler struct {
	medicationRepo *repository.MedicationRepository
}

func NewMedicationHandler(medicationRepo *repository.MedicationRepository) *MedicationHandler {
	return &MedicationHandler{
		medicationRepo: medicationRepo,
	}
}

func (h *MedicationHandler) RegisterLog(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		errors.HandleBadRequest(c, "無効なユーザーIDです", err)
		return
	}

	// リクエストボディを構造体にバインド
	var req dto.MedicationLogRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		errors.HandleValidationError(c, "リクエストボディが無効です", bindErr)
		return
	}

	log.Info().
		Str("user_id", userID).
		Bool("has_bleeding", req.HasBleeding).
		Msg("服用記録の登録を開始します")

	medicationLog := model.MedicationLog{
		HasBleeding: req.HasBleeding,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 日付が指定されている場合は、その日付を使用
	if req.Date != nil {
		medicationLog.CreatedAt = *req.Date
	}

	// リポジトリを呼び出す
	ctx := c.Request.Context()
	err = h.medicationRepo.RegisterLogWithContext(ctx, userID, medicationLog)
	if err != nil {
		errors.HandleDatabaseError(c, "服用記録登録", err)
		return
	}

	log.Info().
		Str("user_id", userID).
		Msg("服用記録の登録が完了しました")

	c.JSON(200, dto.BaseResponse{
		Success: true,
		Message: "medication log registered successfully",
	})
}

// GetLogs はユーザーの服用記録を取得するハンドラー
func (h *MedicationHandler) GetLogs(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		errors.HandleBadRequest(c, "無効なユーザーIDです", err)
		return
	}

	log.Debug().
		Str("user_id", userID).
		Msg("服用記録の取得を開始します")

	// 服用記録を取得
	ctx := c.Request.Context()
	logs, err := h.medicationRepo.GetLogsByUserIDWithContext(ctx, userID)
	if err != nil {
		errors.HandleDatabaseError(c, "服用記録取得", err)
		return
	}

	log.Info().
		Str("user_id", userID).
		Int("count", len(logs)).
		Msg("服用記録の取得が完了しました")

	c.JSON(200, logs)
}

// GetLogByID は特定のIDの服薬ログを取得するハンドラー
func (h *MedicationHandler) GetLogByID(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// URLからIDパラメータを取得
	logIDStr := c.Param("id")
	logID, err := strconv.ParseUint(logIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	// 服薬ログを取得
	log, err := h.medicationRepo.GetLogByID(userID, uint(logID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "medication log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// UpdateLog は指定されたIDの服薬ログを更新するハンドラー
func (h *MedicationHandler) UpdateLog(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// URLからIDパラメータを取得
	logIDStr := c.Param("id")
	logID, err := strconv.ParseUint(logIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	// リクエストボディを構造体にバインド
	var req dto.MedicationLogRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.medicationRepo.UpdateLog(userID, uint(logID), req.HasBleeding)
	if err != nil {
		if err.Error() == "log not found or user not authorized" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update medication log"})
		return

	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		Success: true,
		Message: "medication log updated successfully",
	})
}

// GetMedicationStatus は現在の服薬ステータスを取得するハンドラー
func (h *MedicationHandler) GetMedicationStatus(c *gin.Context) {
	// ユーザーIDを取得
	userID, err := helper.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// サービスから服薬ステータスを取得
	medicationService := service.NewMedicationService(h.medicationRepo)
	status, err := medicationService.GetMedicationStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get medication status"})
		return
	}

	c.JSON(http.StatusOK, status)
}
