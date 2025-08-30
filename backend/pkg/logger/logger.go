package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger はログ設定を初期化する
func InitLogger() {
	// ログレベル設定
	level := getLogLevel()
	zerolog.SetGlobalLevel(level)

	// Lambda環境ではJSON形式、ローカルでは見やすい形式
	if isLambdaEnvironment() {
		// CloudWatch用のJSON形式
		log.Logger = zerolog.New(os.Stderr).With().
			Timestamp().
			Str("service", "okusuri-backend").
			Logger()
	} else {
		// ローカル開発用の見やすい形式
		log.Logger = zerolog.New(
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
				NoColor:    false,
			},
		).With().
			Timestamp().
			Str("service", "okusuri-backend").
			Logger()
	}

	log.Info().Str("level", level.String()).Msg("Logger initialized")
}

// getLogLevel は環境変数からログレベルを取得する
func getLogLevel() zerolog.Level {
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch levelStr {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// isLambdaEnvironment はLambda環境かどうかを判定する
func isLambdaEnvironment() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
