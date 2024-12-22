package rcon

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/cenkalti/backoff/v4"
	"github.com/gorcon/rcon"
	"github.com/nekomeowww/factorio-rcon-api/v2/internal/configs"
	"github.com/nekomeowww/fo"
	"github.com/nekomeowww/xo/logger"
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

//counterfeiter:generate -o fake/rcon.go --fake-name FakeRCON . RCON//counterfeiter:generate -o fake/rcon.go --fake-name FakeRCON . RCON
type RCON interface {
	Close() error
	Execute(ctx context.Context, command string) (string, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	IsReady() bool
}

var _ RCON = (*RCONConn)(nil)

type RCONConn struct {
	*rcon.Conn

	host     string
	port     string
	password string

	ready         atomic.Bool
	reconnectChan chan struct{}
	readyChan     chan struct{}

	mutex  sync.RWMutex
	logger *logger.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRCON() func(NewRCONParams) (RCON, error) {
	return func(params NewRCONParams) (RCON, error) {
		connWrapper := &RCONConn{
			Conn:          nil,
			mutex:         sync.RWMutex{},
			logger:        params.Logger,
			host:          params.Config.Factorio.RCONHost,
			port:          params.Config.Factorio.RCONPort,
			password:      params.Config.Factorio.RCONPassword,
			reconnectChan: make(chan struct{}, 1),
			readyChan:     make(chan struct{}, 1),
		}

		ctx, cancel := context.WithCancel(context.Background())
		connWrapper.ctx = ctx
		connWrapper.cancel = cancel

		// Start the connection manager
		go connWrapper.connectionManager()

		// Trigger initial connection
		select {
		case connWrapper.reconnectChan <- struct{}{}:
		default:
		}

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return fo.Invoke0(ctx, func() error {
					connWrapper.cancel()
					close(connWrapper.reconnectChan)
					close(connWrapper.readyChan)

					connWrapper.mutex.Lock()
					defer connWrapper.mutex.Unlock()

					if connWrapper.Conn != nil {
						return connWrapper.Conn.Close()
					}
					return nil
				})
			},
		})

		return connWrapper, nil
	}
}

func (r *RCONConn) connectionManager() {
	backoffStrategy := backoff.NewExponentialBackOff()

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-r.reconnectChan:
			r.ready.Store(false)

			err := fo.Invoke0(r.ctx, func() error {
				return backoff.Retry(func() error {
					return r.establishConnection(r.ctx)
				}, backoffStrategy)
			})

			if err != nil {
				r.logger.Error("failed to establish RCON connection after retries", zap.Error(err))
				continue
			}

			r.ready.Store(true)

			select {
			case r.readyChan <- struct{}{}:
			default:
			}
		}
	}
}

func (r *RCONConn) establishConnection(ctx context.Context) error {
	return fo.Invoke0(ctx, func() error {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		if r.Conn != nil {
			_ = r.Conn.Close()
		}

		conn, err := rcon.Dial(net.JoinHostPort(r.host, r.port), r.password)
		if err != nil {
			r.logger.Error("failed to connect to RCON", zap.Error(err))
			return err
		}

		r.Conn = conn

		// Test the connection
		_, err = r.Conn.Execute("/help")
		if err != nil {
			r.logger.Error("failed to ping RCON", zap.Error(err))
			return err
		}

		r.logger.Info("RCON connection established successfully")

		return nil
	})
}

func (r *RCONConn) Execute(ctx context.Context, command string) (string, error) {
	return fo.Invoke(ctx, func() (string, error) {
		if !r.IsReady() {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-r.readyChan:
			}
		}

		r.mutex.RLock()
		conn := r.Conn
		r.mutex.RUnlock()
		if conn == nil {
			return r.Execute(ctx, command)
		}

		resp, err := conn.Execute(command)
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") &&
				!strings.Contains(err.Error(), "connection reset by peer") &&
				!errors.Is(err, io.EOF) {
				return "", err
			}

			r.logger.Warn("RCON connection lost, reconnecting...")

			select {
			case r.reconnectChan <- struct{}{}:
			default:
			}

			return r.Execute(ctx, command)
		}

		return resp, nil
	})
}

func (r *RCONConn) IsReady() bool {
	return r.ready.Load()
}
