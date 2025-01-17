package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func StaticWithBytes(specBytes []byte, contentType string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Blob(http.StatusOK, contentType, specBytes)
	}
}
