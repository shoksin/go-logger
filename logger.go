package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GlobalLogger — глобальный экземпляр *zap.Logger, используемый во всех сервисах.
var GlobalLogger *zap.Logger

// Config — структура для настройки логгера.
type Config struct {
	Level      string
	Service    string
	JSONFormat bool //1 - JSON, 0 - text
}

// InitLogger инициализирует GlobalLogger с учётом Config.
// Вызывать нужно единожды в каждом сервисе (обычно в cmd/main.go).
func InitLogger(cfg Config) error {

	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if cfg.JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	// Выводим в stdout (чтобы Docker/Kubernetes мог собирать логи из stdout)
	writer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(encoder, writer, level)

	options := []zap.Option{
		zap.AddCaller(), // добавляет caller (file:line)
		zap.Fields(zap.String("service", cfg.Service)), //"service : serviceName"
	}

	GlobalLogger = zap.New(core, options...)
	return nil
}

// Sync вызывает Sync() у zap.Logger, чтобы убедиться, что все буферы сброшены.
// Рекомендуется defer logger.Sync() в main.go каждого сервиса.
func Sync() {
	if GlobalLogger != nil {
		_ = GlobalLogger.Sync()
	}
}
