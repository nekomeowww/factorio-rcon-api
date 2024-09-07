package rcon

import (
	"context"
	"net"

	"github.com/gorcon/rcon"
	"github.com/nekomeowww/factorio-rcon-api/internal/configs"
	"go.uber.org/fx"
)

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(NewRCON()),
	)
}

type NewRCONParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *configs.Config
}

func NewRCON() func(NewRCONParams) (*rcon.Conn, error) {
	return func(params NewRCONParams) (*rcon.Conn, error) {
		conn, err := rcon.Dial(net.JoinHostPort(params.Config.Factorio.RCONHost, params.Config.Factorio.RCONPort), params.Config.Factorio.RCONPassword)
		if err != nil {
			return nil, err
		}

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(context.Context) error {
				return conn.Close()
			},
		})

		return conn, nil
	}
}
