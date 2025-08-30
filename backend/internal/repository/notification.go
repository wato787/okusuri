package repository

import (
	"context"
	"fmt"
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"
	"time"

	"github.com/guregu/dynamo/v2"
)

type NotificationRepository struct {
	table dynamo.Table
}

func NewNotificationRepository() *NotificationRepository {
	db := config.GetDB()
	table := db.Table(config.GetDynamoDBTableName())

	return &NotificationRepository{
		table: table,
	}
}

// GetSetting はユーザーの通知設定をDynamoDBから取得する
func (r *NotificationRepository) GetSetting(userID, platform string) (*model.NotificationSetting, error) {
	pk := fmt.Sprintf("USER#%s", userID)
	sk := fmt.Sprintf("NOTIFICATION#%s", platform)

	var result model.OkusuriTable
	err := r.table.Get("PK", pk).Range("SK", dynamo.Equal, sk).One(context.Background(), &result)
	if err != nil {
		return nil, err
	}

	// OkusuriTableからNotificationSettingに変換
	setting := &model.NotificationSetting{
		Platform:     getStringValue(result.Data, "platform", platform),
		IsEnabled:    getBoolValue(result.Data, "isEnabled", true),
		Subscription: getStringValue(result.Data, "subscription", ""),
		CreatedAt:    parseTime(getStringValue(result.Data, "createdAt", "")),
		UpdatedAt:    parseTime(getStringValue(result.Data, "updatedAt", "")),
	}

	return setting, nil
}

// RegisterSetting はユーザーの通知設定をDynamoDBに登録/更新する
func (r *NotificationRepository) RegisterSetting(userID string, setting model.NotificationSetting) error {
	pk := fmt.Sprintf("USER#%s", userID)
	sk := fmt.Sprintf("NOTIFICATION#%s", setting.Platform)

	// OkusuriTable形式でデータを保存
	item := model.OkusuriTable{
		PK:   pk,
		SK:   sk,
		Type: "NOTIFICATION",
		Data: map[string]interface{}{
			"platform":     setting.Platform,
			"isEnabled":    setting.IsEnabled,
			"subscription": setting.Subscription,
			"createdAt":    setting.CreatedAt.Format(time.RFC3339),
			"updatedAt":    setting.UpdatedAt.Format(time.RFC3339),
		},
		CreatedAt: setting.CreatedAt.Format(time.RFC3339),
		UpdatedAt: setting.UpdatedAt.Format(time.RFC3339),
	}

	// DynamoDBに保存
	err := r.table.Put(item).Run(context.Background())
	return err
}

// ヘルパー関数はmedication.goで定義済み
