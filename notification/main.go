package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"okusuri-notification/pkg/config"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/guregu/dynamo/v2"
)

// DynamoDBテーブル構造（単一テーブル設計）
type OkusuriTable struct {
	PK        string                 `dynamo:"PK,hash"`                    // Partition Key
	SK        string                 `dynamo:"SK,range"`                   // Sort Key
	Date      string                 `dynamo:"Date,index:DateIndex,hash"`   // GSI1: 日付検索
	Data      map[string]interface{} `dynamo:"Data"`                       // エンティティ固有のデータ
	CreatedAt string                 `dynamo:"CreatedAt"`                   // 作成日時 (ISO8601)
	UpdatedAt string                 `dynamo:"UpdatedAt"`                   // 更新日時 (ISO8601)
	TTL       *int64                 `dynamo:"TTL,omitempty"`               // TTL（必要に応じて）
}

// モデル定義（Cognitoから取得するユーザー情報）
type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
	Image         *string   `json:"image"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// 通知設定（DynamoDBから取得）
type NotificationSetting struct {
	Platform     string `json:"platform"`
	IsEnabled    bool   `json:"isEnabled"`
	Subscription string `json:"subscription"`
}

// 服用履歴（DynamoDBから取得）
type MedicationLog struct {
	HasBleeding bool      `json:"hasBleeding"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// DTOとレスポンス構造体
type MedicationStatusResponse struct {
	CurrentStreak       int  `json:"currentStreak"`
	IsRestPeriod        bool `json:"isRestPeriod"`
	RestDaysLeft        int  `json:"restDaysLeft"`
	ConsecutiveBleeding int  `json:"consecutiveBleeding"`
}

// Push通知関連の構造体
type PushSubscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

type NotificationData struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
}

// リポジトリ層
type Repository struct {
	table         dynamo.Table
	cognitoClient *cognitoidentityprovider.Client
}

func NewRepository(db *dynamo.DB, cognitoClient *cognitoidentityprovider.Client) *Repository {
	return &Repository{
		table:         db.Table(config.GetDynamoDBTableName()),
		cognitoClient: cognitoClient,
	}
}

// Cognitoからユーザー情報を取得
func (r *Repository) GetUsers() ([]User, error) {
	// CognitoのListUsers APIを呼び出し
	input := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: &[]string{config.GetCognitoUserPoolID()}[0],
	}

	result, err := r.cognitoClient.ListUsers(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("Cognitoユーザー取得エラー: %v", err)
	}

	var users []User
	for _, cognitoUser := range result.Users {
		user := unmarshalCognitoUser(cognitoUser)
		users = append(users, *user)
	}

	return users, nil
}

// DynamoDBから通知設定を取得
func (r *Repository) GetNotificationSettings() ([]NotificationSetting, error) {
	var results []OkusuriTable
	err := r.table.Scan().Filter("begins_with(SK, 'NOTIFICATION#')").All(context.Background(), &results)
	if err != nil {
		return nil, fmt.Errorf("通知設定取得エラー: %v", err)
	}

	var settings []NotificationSetting
	for _, result := range results {
		if data, ok := result.Data["platform"].(string); ok {
			setting := NotificationSetting{
				Platform:     data,
				IsEnabled:    getBoolValue(result.Data, "isEnabled", true),
				Subscription: getStringValue(result.Data, "subscription", ""),
			}
			settings = append(settings, setting)
		}
	}

	return settings, nil
}

// DynamoDBから服用履歴を取得
func (r *Repository) GetMedicationLogs(userID string) ([]MedicationLog, error) {
	var results []OkusuriTable
	err := r.table.Get("PK", "USER#"+userID).
		Filter("begins_with(SK, 'MEDICATION#')").
		All(context.Background(), &results)
	if err != nil {
		return nil, fmt.Errorf("服用履歴取得エラー: %v", err)
	}

	var logs []MedicationLog
	for _, result := range results {
		log := MedicationLog{
			HasBleeding: getBoolValue(result.Data, "hasBleeding", false),
			CreatedAt:   parseTime(getStringValue(result.Data, "createdAt", "")),
			UpdatedAt:   parseTime(getStringValue(result.Data, "updatedAt", "")),
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// ヘルパー関数
func unmarshalCognitoUser(cognitoUser types.UserType) *User {
	user := &User{
		ID:        *cognitoUser.Username,
		CreatedAt: *cognitoUser.UserCreateDate,
		UpdatedAt: *cognitoUser.UserLastModifiedDate,
	}

	// 属性から値を取得
	for _, attr := range cognitoUser.Attributes {
		switch *attr.Name {
		case "name":
			user.Name = *attr.Value
		case "email":
			user.Email = *attr.Value
		case "email_verified":
			user.EmailVerified = *attr.Value == "true"
		case "picture":
			user.Image = attr.Value
		}
	}

	return user
}

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

// サービス層
type NotificationService struct {
	recentSends     map[string]time.Time
	recentSendMutex sync.Mutex
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		recentSends: make(map[string]time.Time),
	}
}

func (s *NotificationService) isRecentlySent(subKey string) bool {
	s.recentSendMutex.Lock()
	defer s.recentSendMutex.Unlock()

	lastSent, exists := s.recentSends[subKey]
	if !exists {
		return false
	}

	timeSinceLast := time.Since(lastSent)
	log.Printf("前回の送信からの経過時間: %v (サブスクリプション: %s...)",
		timeSinceLast.Round(time.Second), subKey[:min(10, len(subKey))])
	return timeSinceLast < 5*time.Minute
}

func (s *NotificationService) markAsSent(subKey string) {
	s.recentSendMutex.Lock()
	defer s.recentSendMutex.Unlock()

	s.recentSends[subKey] = time.Now()
	log.Printf("サブスクリプション %s... を送信済みとしてマークしました", subKey[:min(10, len(subKey))])

	// 古い記録をクリーンアップ（1時間以上前のものを削除）
	for key, lastSent := range s.recentSends {
		if time.Since(lastSent) > time.Hour {
			delete(s.recentSends, key)
			log.Printf("古い送信記録を削除: %s...", key[:min(10, len(key))])
		}
	}
}

func (s *NotificationService) SendNotificationWithDays(
	user User, setting NotificationSetting, message string, consecutiveDays int,
) error {
	if setting.Subscription == "" {
		log.Printf("ユーザーID: %s のサブスクリプションが空です", user.ID)
		return fmt.Errorf("サブスクリプションが見つかりません")
	}

	subscriptionPreview := setting.Subscription
	if len(subscriptionPreview) > 10 {
		subscriptionPreview = subscriptionPreview[:10] + "..."
	}

	log.Printf("ユーザーID: %s の処理を開始します", user.ID)
	log.Printf("サブスクリプション: %s", subscriptionPreview)

	var subscription PushSubscription
	err := json.Unmarshal([]byte(setting.Subscription), &subscription)
	if err != nil {
		log.Printf("サブスクリプションのパースに失敗: %v", err)
		return fmt.Errorf("サブスクリプションのパースに失敗: %v", err)
	}

	subKey := subscription.Endpoint
	if s.isRecentlySent(subKey) {
		log.Printf("サブスクリプション %s は最近送信済みのためスキップします", subscriptionPreview)
		return nil
	}

	vapidPublicKey := config.GetVAPIDPublicKey()
	vapidPrivateKey := config.GetVAPIDPrivateKey()

	if vapidPublicKey == "" || vapidPrivateKey == "" {
		log.Printf("VAPID鍵が設定されていません")
		return fmt.Errorf("VAPID鍵が設定されていません")
	}

	notificationData := NotificationData{
		Title: "お薬通知",
		Body:  message,
		Data: map[string]string{
			"messageId":       fmt.Sprintf("medication-%d", time.Now().UnixNano()),
			"timestamp":       fmt.Sprintf("%d", time.Now().Unix()),
			"userId":          user.ID,
			"consecutiveDays": fmt.Sprintf("%d", consecutiveDays),
		},
	}

	payload, err := json.Marshal(notificationData)
	if err != nil {
		log.Printf("通知内容のJSON変換に失敗: %v", err)
		return fmt.Errorf("通知内容のJSON変換に失敗: %v", err)
	}

	_, err = webpush.SendNotification(
		payload,
		&webpush.Subscription{
			Endpoint: subscription.Endpoint,
			Keys: webpush.Keys{
				P256dh: subscription.Keys.P256dh,
				Auth:   subscription.Keys.Auth,
			},
		},
		&webpush.Options{
			VAPIDPublicKey:  vapidPublicKey,
			VAPIDPrivateKey: vapidPrivateKey,
			TTL:             30,
			Subscriber:      "example@example.com",
		},
	)

	if err != nil {
		log.Printf("通知送信エラー: %v", err)
		return fmt.Errorf("通知送信エラー: %v", err)
	}

	s.markAsSent(subKey)
	log.Printf("通知送信成功 - ユーザーID: %s", user.ID)

	return nil
}

// 薬のステータス計算
func calculateMedicationStatus(logs []MedicationLog) *MedicationStatusResponse {
	if len(logs) == 0 {
		return &MedicationStatusResponse{
			CurrentStreak:       0,
			IsRestPeriod:        false,
			RestDaysLeft:        0,
			ConsecutiveBleeding: 0,
		}
	}

	// 最新30日分のログを処理
	today := time.Now().Truncate(24 * time.Hour)
	
	// 最新の服薬から今日まで何日経過したかを計算
	latestLog := logs[0]
	latestDate := latestLog.CreatedAt.Truncate(24 * time.Hour)
	daysSinceLastMedication := int(today.Sub(latestDate).Hours() / 24)

	// 4日以上経過していれば休薬期間
	if daysSinceLastMedication >= 4 {
		restDaysLeft := max(0, 7-daysSinceLastMedication)
		return &MedicationStatusResponse{
			CurrentStreak:       0,
			IsRestPeriod:        true,
			RestDaysLeft:        restDaysLeft,
			ConsecutiveBleeding: 0,
		}
	}

	// 連続服薬日数を計算
	currentStreak := 0
	consecutiveBleeding := 0
	
	// 今日から逆算して連続服薬日数を計算
	checkDate := today
	for i := 0; i < 30; i++ {
		found := false
		var todayLog *MedicationLog
		
		for _, log := range logs {
			logDate := log.CreatedAt.Truncate(24 * time.Hour)
			if logDate.Equal(checkDate) {
				found = true
				todayLog = &log
				break
			}
		}
		
		if !found {
			break
		}
		
		currentStreak++
		
		// 出血が続いている日数を計算
		if todayLog.HasBleeding {
			consecutiveBleeding++
		} else {
			consecutiveBleeding = 0
		}
		
		checkDate = checkDate.AddDate(0, 0, -1)
	}

	return &MedicationStatusResponse{
		CurrentStreak:       currentStreak,
		IsRestPeriod:        false,
		RestDaysLeft:        0,
		ConsecutiveBleeding: consecutiveBleeding,
	}
}

// メッセージ生成
func generateStatusBasedMessage(status *MedicationStatusResponse) string {
	if status.IsRestPeriod {
		if status.RestDaysLeft > 0 {
			return fmt.Sprintf("現在休薬期間中です。あと%d日で服薬を再開してください。", status.RestDaysLeft)
		} else {
			return "休薬期間が終了しました。本日から服薬を再開してください。"
		}
	} else {
		if status.CurrentStreak > 0 {
			return fmt.Sprintf("お薬の時間です。忘れずに服用してください。（連続%d日目）", status.CurrentStreak)
		} else {
			return "お薬の時間です。忘れずに服用してください。"
		}
	}
}

// メイン処理
func handleRequest(ctx context.Context, event interface{}) (interface{}, error) {
	requestTime := time.Now()
	log.Printf("========== 通知送信処理開始 [%s] ==========", requestTime.Format("2006-01-02 15:04:05"))

	// AWS設定
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("AWS設定エラー: %v", err)
		return nil, fmt.Errorf("AWS設定エラー: %v", err)
	}

	// DynamoDB接続
	db := dynamo.New(cfg)
	
	// Cognitoクライアント
	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)

	// リポジトリ初期化
	repo := NewRepository(db, cognitoClient)

	// ユーザー一覧を取得（Cognitoから）
	users, err := repo.GetUsers()
	if err != nil {
		log.Printf("ユーザー取得エラー: %v", err)
		return nil, fmt.Errorf("ユーザー取得エラー: %v", err)
	}
	log.Printf("取得したユーザー数: %d", len(users))

	// 通知設定一覧を取得（DynamoDBから）
	settings, err := repo.GetNotificationSettings()
	if err != nil {
		log.Printf("通知設定取得エラー: %v", err)
		return nil, fmt.Errorf("通知設定取得エラー: %v", err)
	}
	log.Printf("取得した通知設定数: %d", len(settings))

	// 通知設定をユーザーIDでマップ化
	settingsMap := make(map[string]NotificationSetting)
	for _, setting := range settings {
		// ユーザーIDは通知設定から取得する必要がある
		// ここでは簡略化のため、設定の数だけ処理
		settingsMap[fmt.Sprintf("user_%d", len(settingsMap))] = setting
	}
	log.Printf("通知対象ユーザー数: %d", len(settingsMap))

	// 通知サービスを初期化
	notificationSvc := NewNotificationService()
	sentSubs := make(map[string]bool)
	sentCount := 0

	log.Println("----- 通知送信処理開始 -----")

	for _, user := range users {
		// ユーザーに対応する通知設定を取得
		setting, ok := settingsMap[user.ID]
		if !ok || !setting.IsEnabled {
			continue
		}

		if _, alreadySent := sentSubs[setting.Subscription]; alreadySent && setting.Subscription != "" {
			continue
		}

		// 薬のステータスを取得してメッセージを生成
		message := "お薬の時間です。忘れずに服用してください。"
		consecutiveDays := 0
		
		medicationLogs, statusErr := repo.GetMedicationLogs(user.ID)
		if statusErr == nil {
			medicationStatus := calculateMedicationStatus(medicationLogs)
			message = generateStatusBasedMessage(medicationStatus)
			consecutiveDays = medicationStatus.CurrentStreak
		}

		sendErr := notificationSvc.SendNotificationWithDays(user, setting, message, consecutiveDays)
		if sendErr != nil {
			log.Printf("通知送信失敗: %v", sendErr)
			continue
		}

		if setting.Subscription != "" {
			sentSubs[setting.Subscription] = true
		}
		sentCount++
		log.Printf("ユーザーID: %s への通知送信成功", user.ID)
	}

	processingTime := time.Since(requestTime)
	log.Printf("----- 通知送信処理完了: 合計%d件送信 -----", sentCount)
	log.Printf("処理時間: %v", processingTime)
	log.Printf("========== 通知送信処理終了 [%s] ==========\n", time.Now().Format("2006-01-02 15:04:05"))

	return map[string]interface{}{
		"message":         "notification sent successfully",
		"sent_count":      sentCount,
		"process_time_ms": processingTime.Milliseconds(),
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	lambda.Start(handleRequest)
}