package logger

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if env == "prod" {
		cfg := zap.NewProductionConfig()
		logger, err = cfg.Build()
		if err != nil {
			return nil, errors.Wrap(err, "setupLogger")
		}
	} else {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder := zapcore.NewConsoleEncoder(encoderConfig)

		consoleOutput := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(encoder, consoleOutput, zapcore.DebugLevel)

		logger = zap.New(core)
	}

	return logger, nil
}
