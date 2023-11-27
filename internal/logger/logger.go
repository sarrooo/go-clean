package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevels = map[string]zapcore.Level{
	"debug":  zap.DebugLevel,
	"info":   zap.InfoLevel,
	"warn":   zap.WarnLevel,
	"error":  zap.ErrorLevel,
	"dpanic": zap.DPanicLevel,
	"panic":  zap.PanicLevel,
	"fatal":  zap.FatalLevel,
}

func New() (*zap.Logger, error) {
	logLevel := viper.GetString("LOG_LEVEL")
	atom := zap.NewAtomicLevel()
	atom.SetLevel(logLevels[logLevel])
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeName = zapcore.FullNameEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.Level = atom
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
