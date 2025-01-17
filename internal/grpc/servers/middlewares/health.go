package middlewares

import (
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/labstack/echo/v4"
)

func HealthCheck(checkOptions ...health.CheckerOption) echo.HandlerFunc {
	return func(c echo.Context) error {
		opts := make([]health.CheckerOption, 0)
		opts = append(opts,
			health.WithCacheDuration(time.Second),
			health.WithTimeout(time.Second*10), //nolint:mnd
		)

		checker := health.NewChecker(opts...)
		handler := health.NewHandler(checker)

		handler.ServeHTTP(c.Response().Writer, c.Request())

		return nil
	}
}
