package model

import (
	"time"
)

// MedicationLog は服用履歴の構造体（DynamoDB対応）
type MedicationLog struct {
	HasBleeding bool      `json:"hasBleeding"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NotificationSetting は通知設定の構造体（DynamoDB対応）
type NotificationSetting struct {
	Platform     string    `json:"platform"`
	IsEnabled    bool      `json:"isEnabled"`
	Subscription string    `json:"subscription"` // Web Push用のサブスクリプション
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// OkusuriTable はDynamoDB単一テーブル設計のメイン構造体
type OkusuriTable struct {
	PK        string                 `dynamo:"PK"`               // Partition Key
	SK        string                 `dynamo:"SK"`               // Sort Key
	GSI1PK    string                 `dynamo:"GSI1PK,omitempty"` // GSI1 Partition Key
	GSI1SK    string                 `dynamo:"GSI1SK,omitempty"` // GSI1 Sort Key
	Type      string                 `dynamo:"Type,omitempty"`   // レコードタイプ
	Date      string                 `dynamo:"Date,omitempty"`   // 日付（YYYY-MM-DD形式）
	Data      map[string]interface{} `dynamo:"Data,omitempty"`   // データペイロード
	CreatedAt string                 `dynamo:"CreatedAt"`        // 作成日時（ISO8601）
	UpdatedAt string                 `dynamo:"UpdatedAt"`        // 更新日時（ISO8601）
	TTL       int64                  `dynamo:"TTL,omitempty"`    // TTL（必要に応じて）
}

// TableName はDynamoDBのテーブル名を返す
func (o OkusuriTable) TableName() string {
	return "okusuri-table"
}
