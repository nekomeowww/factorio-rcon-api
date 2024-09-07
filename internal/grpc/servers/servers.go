package servers

import (
	"github.com/nekomeowww/factorio-rcon-api/internal/grpc/servers/factorioapi/v1/apiserver"
	"go.uber.org/fx"
)

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(apiserver.NewGRPCServer()),
		fx.Provide(apiserver.NewGatewayServer()),
	)
}
