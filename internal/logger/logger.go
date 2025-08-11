package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger // переменная для использования логгера

// Init инициализирует логгер в зависимости от режима логирования
func Init(mode string, logFilePath string) error {
	var (
		err      error
		encoder  zapcore.Encoder
		logLevel zapcore.Level
	)

	if logFilePath == "" {
		return fmt.Errorf("log file path is not set")
	}

	logFilePath = filepath.Clean(logFilePath)

	// Создаём директорию для логов, если её нет
	logDir := filepath.Dir(logFilePath)
	err = os.MkdirAll(logDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	switch mode {
	case "dev":
		logLevel = zap.DebugLevel
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	case "prod":
		logLevel = zap.InfoLevel
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "ts"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	default:
		return fmt.Errorf("unknown SERVER_LOGGING mode: %s", mode)
	}

	consoleCore := zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stdout),
		logLevel,
	)

	fileCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(file),
		logLevel,
	)

	core := zapcore.NewTee(consoleCore, fileCore)

	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Sync очищает все буферизованные записи журнала
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}
