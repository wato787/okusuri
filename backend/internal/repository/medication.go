package repository

import (
	"context"
	"fmt"
	"okusuri-backend/internal/model"
	"okusuri-backend/pkg/config"
	"time"

	"github.com/guregu/dynamo/v2"
)

type MedicationRepository struct {
	table dynamo.Table
}

func NewMedicationRepository() *MedicationRepository {
	db := config.GetDB()
	table := db.Table("okusuri-table")

	return &MedicationRepository{
		table: table,
	}
}

// RegisterLog はユーザーの服用記録をDynamoDBに登録する（後方互換性）
func (r *MedicationRepository) RegisterLog(userID string, log model.MedicationLog) error {
	return r.RegisterLogWithContext(context.Background(), userID, log)
}

// RegisterLogWithContext はユーザーの服用記録をDynamoDBに登録する
func (r *MedicationRepository) RegisterLogWithContext(ctx context.Context, userID string, log model.MedicationLog) error {
	// DynamoDBの単一テーブル設計に基づくキー生成
	pk := fmt.Sprintf("USER#%s", userID)
	sk := fmt.Sprintf("MEDICATION#%s#%d", log.CreatedAt.Format("2006-01-02"), time.Now().UnixNano())

	// OkusuriTable形式でデータを保存
	item := model.OkusuriTable{
		PK:   pk,
		SK:   sk,
		Date: log.CreatedAt.Format("2006-01-02"),
		Data: map[string]interface{}{
			"hasBleeding": log.HasBleeding,
			"createdAt":   log.CreatedAt.Format(time.RFC3339),
			"updatedAt":   log.UpdatedAt.Format(time.RFC3339),
		},
		CreatedAt: log.CreatedAt.Format(time.RFC3339),
		UpdatedAt: log.UpdatedAt.Format(time.RFC3339),
	}

	// DynamoDBに保存
	err := r.table.Put(item).Run(ctx)
	return err
}

// GetLogsByUserID はユーザーIDに基づいて服用履歴をDynamoDBから取得する（後方互換性）
func (r *MedicationRepository) GetLogsByUserID(userID string) ([]model.MedicationLog, error) {
	return r.GetLogsByUserIDWithContext(context.Background(), userID)
}

// GetLogsByUserIDWithContext はユーザーIDに基づいて服用履歴をDynamoDBから取得する
func (r *MedicationRepository) GetLogsByUserIDWithContext(ctx context.Context, userID string) ([]model.MedicationLog, error) {
	pk := fmt.Sprintf("USER#%s", userID)

	var results []model.OkusuriTable
	err := r.table.Get("PK", pk).
		Filter("begins_with(SK, ?)", "MEDICATION#").
		All(ctx, &results)

	if err != nil {
		return nil, err
	}

	// OkusuriTableからMedicationLogに変換
	var logs []model.MedicationLog
	for _, result := range results {
		log := model.MedicationLog{
			HasBleeding: getBoolValue(result.Data, "hasBleeding", false),
			CreatedAt:   parseTime(getStringValue(result.Data, "createdAt", "")),
			UpdatedAt:   parseTime(getStringValue(result.Data, "updatedAt", "")),
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetLogByID はIDに基づいて単一の服薬ログを取得する
func (r *MedicationRepository) GetLogByID(userID string, logID uint) (*model.MedicationLog, error) {
	// DynamoDBでは直接的なID検索は困難なため、ユーザーの全ログから検索
	logs, err := r.GetLogsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 簡易的なID検索（実際の運用では別の方法を検討）
	for _, log := range logs {
		// TODO: より効率的なID検索の実装
		if log.CreatedAt.Unix() == int64(logID) {
			return &log, nil
		}
	}

	return nil, fmt.Errorf("log not found")
}

// UpdateLog は指定されたIDの服薬ログを更新する
func (r *MedicationRepository) UpdateLog(userID string, logID uint, hasBleeding bool) error {
	// DynamoDBでの更新操作
	pk := fmt.Sprintf("USER#%s", userID)
	sk := fmt.Sprintf("MEDICATION#%s#%d", time.Now().Format("2006-01-02"), logID)

	// DynamoDBのUpdate操作
	err := r.table.Update("PK", pk).
		Range("SK", sk).
		Set("Data.hasBleeding", hasBleeding).
		Set("Data.updatedAt", time.Now().Format(time.RFC3339)).
		Set("UpdatedAt", time.Now().Format(time.RFC3339)).
		Run(context.Background())

	return err
}

// GetConsecutiveDays はユーザーの連続服薬日数を計算する
func (r *MedicationRepository) GetConsecutiveDays(userID string) (int, error) {
	logs, err := r.GetLogsByUserID(userID)
	if err != nil {
		return 0, err
	}

	if len(logs) == 0 {
		return 0, nil
	}

	// 連続日数をカウント
	consecutiveDays := 1
	today := time.Now().Truncate(24 * time.Hour)

	// 最新の記録が今日かどうかチェック
	latestLog := logs[0]
	if !latestLog.CreatedAt.Truncate(24 * time.Hour).Equal(today) {
		return 0, nil
	}

	// 連続日数を計算
	for i := 1; i < len(logs); i++ {
		currentLog := logs[i]
		previousLog := logs[i-1]

		// 日付の差が1日かチェック
		diff := currentLog.CreatedAt.Truncate(24 * time.Hour).Sub(previousLog.CreatedAt.Truncate(24 * time.Hour))
		if diff == 24*time.Hour {
			consecutiveDays++
		} else {
			break
		}
	}

	return consecutiveDays, nil
}

// ヘルパー関数
func getBoolValue(data map[string]interface{}, key string, defaultValue bool) bool {
	if value, ok := data[key].(bool); ok {
		return value
	}
	return defaultValue
}

func getStringValue(data map[string]interface{}, key string, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return defaultValue
}

func parseTime(timeStr string) time.Time {
	if timeStr == "" {
		return time.Now()
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Now()
	}
	return t
}
