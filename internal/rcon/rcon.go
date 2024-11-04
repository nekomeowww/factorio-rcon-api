package rcon

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"errors"
	"io"
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

//counterfeiter:generate -o fake/rcon.go --fake-name FakeRCON . RCON
type RCON interface {
	Close() error
	Execute(ctx context.Context, command string) (string, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	IsReady() bool
}

type RCONConn struct {
	*rcon.Conn

	host     string
	port     string
	password string

	ready          bool
	readinessMutex sync.RWMutex
	mutex          sync.RWMutex
	logger         *logger.Logger
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewRCON() func(NewRCONParams) (RCON, error) {
	return func(params NewRCONParams) (RCON, error) {
		connWrapper := &RCONConn{
			Conn:     nil,
			mutex:    sync.RWMutex{},
			logger:   params.Logger,
			host:     params.Config.Factorio.RCONHost,
			port:     params.Config.Factorio.RCONPort,
			password: params.Config.Factorio.RCONPassword,
		}

		ctx, cancel := context.WithCancel(context.Background())
		connWrapper.ctx = ctx
		connWrapper.cancel = cancel

		go connWrapper.Connect(ctx)

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(context.Context) error {
				connWrapper.cancel()

				connWrapper.mutex.RLock()
				defer connWrapper.mutex.RUnlock()

				if connWrapper.Conn == nil {
					return nil
				}

				_ = connWrapper.Conn.Close()

				return nil
			},
		})

		return connWrapper, nil
	}
}

func (r *RCONConn) connect() (*rcon.Conn, error) {
	conn, err := rcon.Dial(net.JoinHostPort(r.host, r.port), r.password)
	if err != nil {
		r.logger.Error("failed to connect to RCON, will attempt to reconnect", zap.Error(err))

		return nil, err
	}

	return conn, nil
}

func (r *RCONConn) Connect(ctx context.Context) {
	r.setUnready()

	err := fo.Invoke0(ctx, func() error {
		return backoff.Retry(func() error {
			conn, err := r.connect()
			if err != nil {
				return err
			}

			r.mutex.Lock()
			defer r.mutex.Unlock()

			r.Conn = conn

			err = r.ping(ctx)
			if err != nil {
				r.logger.Error("failed to ping RCON, will attempt to reconnect", zap.Error(err))

				return err
			}

			return nil
		}, backoff.NewExponentialBackOff())
	})
	if err != nil {
		r.logger.Error("failed to connect to RCON", zap.Error(err))
		return
	}

	r.setReady()
}

func (r *RCONConn) Execute(ctx context.Context, command string) (string, error) {
	return fo.Invoke(ctx, func() (string, error) {
		_, _, err := lo.AttemptWithDelay(40, 250*time.Millisecond, func(_ int, _ time.Duration) error {
			if !r.IsReady() {
				return ErrTimeout
			}

			return nil
		})
		if err != nil {
			return "", err
		}

		resp, err := r.Conn.Execute(command)
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.logger.Warn("RCON connection is closed, attempting to reconnect")

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				r.Connect(ctx)
				return r.Execute(ctx, command)
			}

			return "", err
		}

		return resp, nil
	})
}

func (r *RCONConn) setUnready() {
	r.readinessMutex.Lock()
	defer r.readinessMutex.Unlock()

	r.ready = false
}

func (r *RCONConn) setReady() {
	r.readinessMutex.Lock()
	defer r.readinessMutex.Unlock()

	r.ready = true
}

func (r *RCONConn) ping(ctx context.Context) error {
	return fo.Invoke0(ctx, func() error {
		_, err := r.Conn.Execute("/help")
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *RCONConn) IsReady() bool {
	r.readinessMutex.RLock()
	r.mutex.RLock()

	defer r.mutex.RUnlock()
	defer r.readinessMutex.RUnlock()

	return r.ready && r.Conn != nil
}
