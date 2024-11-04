package middlewares

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/labstack/echo/v4"
	v1 "github.com/nekomeowww/factorio-rcon-api/apis/factorioapi/v1"
	"github.com/nekomeowww/factorio-rcon-api/pkg/apierrors"
)

func ScalarDocumentation(title string) echo.HandlerFunc {
	return func(c echo.Context) error {
		content, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: string(v1.OpenAPIV3SpecYaml()),
			CustomOptions: scalar.CustomOptions{
				PageTitle: title,
			},
		})
		if err != nil {
			return apierrors.NewErrInternal().WithError(err).WithDetailf("failed to generate API documentation: %s", err.Error()).AsEchoResponse(c)
		}

		return c.HTML(http.StatusOK, content)
	}
}
