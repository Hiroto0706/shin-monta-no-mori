package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var config zap.Config

// init関数は、このパッケージが初めてインポートされたときに自動的に実行されます。
func init() {
	// FIXME: Asia/Tokyo は使えないのでUTCにしている
	time.Local = time.UTC

	// zapの開発環境用のデフォルト設定を取得
	config = zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	// JSONエンコーダーの設定を取得します。これにより、ログメッセージがJSON形式で出力されます。
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig
}
