package middlewares

import "github.com/labstack/echo/v4"

func StaticWithBytes(specBytes []byte, contentType string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Blob(200, contentType, specBytes)
	}
}
