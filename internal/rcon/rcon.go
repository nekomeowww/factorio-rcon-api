package rcon

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorcon/rcon"
	"github.com/nekomeowww/factorio-rcon-api/internal/configs"
	"github.com/nekomeowww/fo"
	"github.com/nekomeowww/xo/logger"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	ErrTimeout = errors.New("RCON connection is not established within deadline threshold")
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
	Logger    *logger.Logger
}

type RCON struct {
	*rcon.Conn

	mutex sync.RWMutex
}

func NewRCON() func(NewRCONParams) (*RCON, error) {
	return func(params NewRCONParams) (*RCON, error) {
		var err error
		var conn *rcon.Conn

		connWrapper := &RCON{
			Conn:  nil,
			mutex: sync.RWMutex{},
		}

		ctx, cancel := context.WithCancel(context.Background())

		go fo.Invoke0(ctx, func() error {
			return backoff.Retry(func() error {
				conn, err = rcon.Dial(net.JoinHostPort(params.Config.Factorio.RCONHost, params.Config.Factorio.RCONPort), params.Config.Factorio.RCONPassword)
				if err != nil {
					params.Logger.Error("failed to connect to RCON, will attempt to reconnect", zap.Error(err))

					return err
				}

				connWrapper.mutex.Lock()
				defer connWrapper.mutex.Unlock()

				connWrapper.Conn = conn

				return nil
			}, backoff.NewExponentialBackOff())
		})

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(context.Context) error {
				cancel()

				connWrapper.mutex.RLock()
				defer connWrapper.mutex.RUnlock()

				if connWrapper.Conn == nil {
					return nil
				}

				return connWrapper.Conn.Close()
			},
		})

		return connWrapper, nil
	}
}

func (r *RCON) Execute(ctx context.Context, command string) (string, error) {
	return fo.Invoke(ctx, func() (string, error) {
		r.mutex.RLock()
		defer r.mutex.RUnlock()

		_, _, err := lo.AttemptWithDelay(10, time.Second, func(_ int, _ time.Duration) error {
			if r.Conn == nil {
				return ErrTimeout
			}

			return nil
		})
		if err != nil {
			return "", err
		}

		return r.Conn.Execute(command)
	})
}
