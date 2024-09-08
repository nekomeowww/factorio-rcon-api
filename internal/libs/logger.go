package libs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nekomeowww/xo"
	"github.com/nekomeowww/xo/logger"
	"github.com/nekomeowww/xo/logger/loki"
	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

func NewLogger() func() (*logger.Logger, error) {
	return func() (*logger.Logger, error) {
		logLevel, err := logger.ReadLogLevelFromEnv()
		if err != nil {
			logLevel = zapcore.InfoLevel
		}

		var isFatalLevel bool
		if logLevel == zapcore.FatalLevel {
			isFatalLevel = true
			logLevel = zapcore.InfoLevel
		}

		logFormat, readFormatError := logger.ReadLogFormatFromEnv()

		logger, err := logger.NewLogger(
			logger.WithLevel(logLevel),
			logger.WithAppName("factorio-rcon-api"),
			logger.WithNamespace("nekomeowww"),
			logger.WithLogFilePath(xo.RelativePathBasedOnPwdOf(filepath.Join("logs", "logs.log"))),
			logger.WithFormat(logFormat),
			logger.WithLokiRemoteConfig(lo.Ternary(os.Getenv("LOG_LOKI_REMOTE_URL") != "", &loki.Config{
				Url:          os.Getenv("LOG_LOKI_REMOTE_URL"),
				BatchMaxSize: 2000,
				BatchMaxWait: 10 * time.Second,
				PrintErrors:  true,
				Labels:       map[string]string{},
			}, nil)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
		if isFatalLevel {
			logger.Error("fatal log level is unacceptable, fallbacks to info level")
		}
		if readFormatError != nil {
			logger.Error("failed to read log format from env, fallbacks to json")
		}

		logger = logger.WithAndSkip(
			1,
		)

		return logger, nil
	}
}
