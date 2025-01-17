package factorioapi

import (
	factorioapiv1 "github.com/nekomeowww/factorio-rcon-api/v2/apis/factorioapi/v1"
	factorioapiv2 "github.com/nekomeowww/factorio-rcon-api/v2/apis/factorioapi/v2"
	consolev1 "github.com/nekomeowww/factorio-rcon-api/v2/internal/grpc/services/factorioapi/v1/console"
	consolev2 "github.com/nekomeowww/factorio-rcon-api/v2/internal/grpc/services/factorioapi/v2/console"
	grpcpkg "github.com/nekomeowww/factorio-rcon-api/v2/pkg/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc/reflection"
)

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(NewFactorioAPI()),
		fx.Provide(consolev1.NewConsoleService()),
		fx.Provide(consolev2.NewConsoleService()),
	)
}

type NewFactorioAPIParams struct {
	fx.In

	Console   *consolev1.ConsoleService
	ConsoleV2 *consolev2.ConsoleService
}

type FactorioAPI struct {
	params *NewFactorioAPIParams
}

func NewFactorioAPI() func(params NewFactorioAPIParams) *FactorioAPI {
	return func(params NewFactorioAPIParams) *FactorioAPI {
		return &FactorioAPI{params: &params}
	}
}

func (c *FactorioAPI) Register(r *grpcpkg.Register) {
	r.RegisterHTTPHandlers([]grpcpkg.HTTPHandler{
		factorioapiv1.RegisterConsoleServiceHandler,
		factorioapiv2.RegisterConsoleServiceHandler,
	})

	r.RegisterGrpcService(func(s reflection.GRPCServer) {
		factorioapiv1.RegisterConsoleServiceServer(s, c.params.Console)
		factorioapiv2.RegisterConsoleServiceServer(s, c.params.ConsoleV2)
	})
}
