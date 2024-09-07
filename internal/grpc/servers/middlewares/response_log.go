package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nekomeowww/xo/logger"
	"go.uber.org/zap"
)

func ResponseLog(logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				return err
			}

			end := time.Now()

			logger.Debug("",
				zap.String("latency", end.Sub(start).String()),
				zap.String("path", c.Request().RequestURI),
				zap.String("remote", c.Request().RemoteAddr),
				zap.String("hosts", c.Request().URL.Host),
				zap.Int("status", c.Response().Status),
				zap.String("method", c.Request().Method),
			)

			return nil
		}
	}
}
