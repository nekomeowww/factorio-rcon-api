package servers

import (
	"github.com/nekomeowww/factorio-rcon-api/v2/internal/grpc/servers/factorioapi/apiserver"
	"go.uber.org/fx"
)

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(apiserver.NewGRPCServer()),
		fx.Provide(apiserver.NewGatewayServer()),
	)
}
