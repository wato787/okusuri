package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// APIError は統一的なエラーレスポンス構造体
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ErrorCode はエラーコードの定数
const (
	ErrCodeInvalidRequest     = "INVALID_REQUEST"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeConflict           = "CONFLICT"
	ErrCodeInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError      = "DATABASE_ERROR"
	ErrCodeValidationFailed   = "VALIDATION_FAILED"
	ErrCodeInvalidUserID      = "INVALID_USER_ID"
	ErrCodeMedicationNotFound = "MEDICATION_NOT_FOUND"
	ErrCodeNotificationError  = "NOTIFICATION_ERROR"
)

// HandleError は統一されたエラーレスポンスを返す
func HandleError(c *gin.Context, statusCode int, errorCode, message string, err error, details ...string) {
	// ログ出力
	logger := log.With().
		Str("path", c.Request.URL.Path).
		Str("method", c.Request.Method).
		Str("user_agent", c.Request.UserAgent()).
		Str("error_code", errorCode).
		Logger()

	if err != nil {
		logger.Error().Err(err).Msg(message)
	} else {
		logger.Warn().Msg(message)
	}

	// レスポンス作成
	apiError := APIError{
		Code:    errorCode,
		Message: message,
	}

	if len(details) > 0 {
		apiError.Details = details[0]
	}

	c.JSON(statusCode, apiError)
}

// HandleBadRequest は400エラーを処理する
func HandleBadRequest(c *gin.Context, message string, err error, details ...string) {
	HandleError(c, http.StatusBadRequest, ErrCodeInvalidRequest, message, err, details...)
}

// HandleUnauthorized は401エラーを処理する
func HandleUnauthorized(c *gin.Context, message string, err error, details ...string) {
	HandleError(c, http.StatusUnauthorized, ErrCodeUnauthorized, message, err, details...)
}

// HandleForbidden は403エラーを処理する
func HandleForbidden(c *gin.Context, message string, err error, details ...string) {
	HandleError(c, http.StatusForbidden, ErrCodeForbidden, message, err, details...)
}

// HandleNotFound は404エラーを処理する
func HandleNotFound(c *gin.Context, message string, err error, details ...string) {
	HandleError(c, http.StatusNotFound, ErrCodeNotFound, message, err, details...)
}

// HandleInternalServerError は500エラーを処理する
func HandleInternalServerError(c *gin.Context, message string, err error, details ...string) {
	HandleError(c, http.StatusInternalServerError, ErrCodeInternalServer, message, err, details...)
}

// HandleDatabaseError はデータベースエラーを処理する
func HandleDatabaseError(c *gin.Context, operation string, err error) {
	message := "データベース操作に失敗しました"
	HandleError(c, http.StatusInternalServerError, ErrCodeDatabaseError, message, err, operation)
}

// HandleValidationError はバリデーションエラーを処理する
func HandleValidationError(c *gin.Context, message string, err error) {
	HandleError(c, http.StatusBadRequest, ErrCodeValidationFailed, message, err)
}
