package services

import (
	"go.uber.org/fx"
	"google.golang.org/grpc/reflection"

	"github.com/nekomeowww/factorio-rcon-api/internal/grpc/services/factorioapi"
	grpcpkg "github.com/nekomeowww/factorio-rcon-api/pkg/grpc"
)

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(NewRegister()),
		fx.Options(factorioapi.Modules()),
	)
}

type NewRegisterParams struct {
	fx.In

	FactorioAPI *factorioapi.FactorioAPI
}

func NewRegister() func(params NewRegisterParams) *grpcpkg.Register {
	return func(params NewRegisterParams) *grpcpkg.Register {
		register := grpcpkg.NewRegister()

		params.FactorioAPI.Register(register)

		register.RegisterGrpcService(func(s reflection.GRPCServer) {
			reflection.Register(s)
		})

		return register
	}
}
