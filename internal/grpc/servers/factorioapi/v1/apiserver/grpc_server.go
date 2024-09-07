package apiserver

import (
	"context"
	"net"

	"github.com/nekomeowww/xo/logger"
	grpcotel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/nekomeowww/factorio-rcon-api/internal/configs"
	"github.com/nekomeowww/factorio-rcon-api/internal/grpc/servers/interceptors"
	grpcpkg "github.com/nekomeowww/factorio-rcon-api/pkg/grpc"
)

type NewGRPCServerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *configs.Config
	Logger    *logger.Logger
	Register  *grpcpkg.Register
}

type GRPCServer struct {
	ListenAddr string

	server   *grpc.Server
	register *grpcpkg.Register
}

func NewGRPCServer() func(params NewGRPCServerParams) *GRPCServer {
	return func(params NewGRPCServerParams) *GRPCServer {
		server := grpc.NewServer(
			grpc.StatsHandler(grpcotel.NewServerHandler()),
			grpc.ChainUnaryInterceptor(
				interceptors.PanicInterceptor(params.Logger),
			),
		)

		params.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				params.Logger.Info("gracefully shutting down gRPC server...")
				server.GracefulStop()
				return nil
			},
		})

		return &GRPCServer{
			ListenAddr: params.Config.APIServer.GrpcServerAddr,
			server:     server,
			register:   params.Register,
		}
	}
}

func RunGRPCServer() func(logger *logger.Logger, server *GRPCServer) error {
	return func(logger *logger.Logger, server *GRPCServer) error {
		for _, serviceRegister := range server.register.GrpcServices {
			serviceRegister(server.server)
		}

		l, err := net.Listen("tcp", server.ListenAddr)
		if err != nil {
			return err
		}

		go func() {
			err := server.server.Serve(l)
			if err != nil && err != grpc.ErrServerStopped {
				logger.Fatal("failed to serve gRPC server", zap.Error(err))
			}
		}()

		logger.Info("gRPC server started", zap.String("listen_addr", server.ListenAddr))

		return nil
	}
}
