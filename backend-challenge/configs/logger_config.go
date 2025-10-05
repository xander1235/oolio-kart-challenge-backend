package configs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger = zap.NewExample(zap.AddCaller())

// InitLogger initializes the logger
func InitLogger(logLevel string) *zap.Logger {
	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)

	logger := zap.New(core, zap.AddCaller()).Named(AppName)
	Logger = logger
	return logger
}
