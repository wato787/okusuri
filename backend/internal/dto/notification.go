package dto

// NotificationSettingRequest は通知設定のリクエスト用DTO
type NotificationSettingRequest struct {
	Platform     string `json:"platform" binding:"required"` // "web", "mobile"等
	IsEnabled    bool   `json:"isEnabled"`                   // 通知の有効/無効
	Subscription string `json:"subscription,omitempty"`      // WebPush用のサブスクリプション
}

// NotificationSettingResponse は通知設定のレスポンス用DTO
type NotificationSettingResponse struct {
	Platform     string `json:"platform"`
	IsEnabled    bool   `json:"isEnabled"`
	Subscription string `json:"subscription,omitempty"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
