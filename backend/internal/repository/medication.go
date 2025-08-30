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
	table := db.Table(config.GetDynamoDBTableName())

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
	// DynamoDBでは直接的な更新が困難なため、削除して再作成
	// 実際の運用では条件付き更新を使用することを推奨
	logs, err := r.GetLogsByUserID(userID)
	if err != nil {
		return err
	}

	for _, log := range logs {
		if log.CreatedAt.Unix() == int64(logID) {
			// 更新されたログを作成
			updatedLog := log
			updatedLog.HasBleeding = hasBleeding
			updatedLog.UpdatedAt = time.Now()

			// 古いログを削除して新しいログを作成
			// TODO: より効率的な更新処理の実装
			return r.RegisterLog(userID, updatedLog)
		}
	}

	return fmt.Errorf("log not found")
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
	latestDate := latestLog.CreatedAt.Truncate(24 * time.Hour)

	// 最新の記録が今日でない場合は連続日数は0
	if !latestDate.Equal(today) {
		return 0, nil
	}

	// 前日から遡って連続日数をカウント
	expectedDate := today.AddDate(0, 0, -1)

	for i := 1; i < len(logs); i++ {
		logDate := logs[i].CreatedAt.Truncate(24 * time.Hour)

		if logDate.Equal(expectedDate) {
			consecutiveDays++
			expectedDate = expectedDate.AddDate(0, 0, -1)
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
