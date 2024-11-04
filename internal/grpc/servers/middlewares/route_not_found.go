package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nekomeowww/factorio-rcon-api/v2/pkg/apierrors"
)

func NotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, apierrors.NewErrNotFound().AsResponse())
}
