package apiserver

import (
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/nekomeowww/xo/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nekomeowww/factorio-rcon-api/internal/configs"
	"github.com/nekomeowww/factorio-rcon-api/internal/grpc/servers/interceptors"
	"github.com/nekomeowww/factorio-rcon-api/internal/grpc/servers/middlewares"
	"github.com/nekomeowww/factorio-rcon-api/internal/libs"
	grpcpkg "github.com/nekomeowww/factorio-rcon-api/pkg/grpc"
	httppkg "github.com/nekomeowww/factorio-rcon-api/pkg/http"
)

type NewGatewayServerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *configs.Config
	Register  *grpcpkg.Register
	Logger    *logger.Logger
	Otel      *libs.Otel
}

type GatewayServer struct {
	ListenAddr     string
	GRPCServerAddr string

	echo   *echo.Echo
	server *http.Server
}

func NewGatewayServer() func(params NewGatewayServerParams) (*GatewayServer, error) {
	return func(params NewGatewayServerParams) (*GatewayServer, error) {
		gob.Register(map[interface{}]interface{}{})

		e := echo.New()
		e.RouteNotFound("/*", middlewares.NotFound)

		e.GET("/apis/docs", middlewares.ScalarDocumentation("Factorio RCON API", "/swagger.json"))

		for path, methodHandlers := range params.Register.EchoHandlers {
			for method, handler := range methodHandlers {
				e.Add(method, path, handler)
			}
		}

		server := &GatewayServer{
			ListenAddr:     params.Config.APIServer.HttpServerAddr,
			GRPCServerAddr: params.Config.APIServer.GrpcServerAddr,
			echo:           e,
			server: &http.Server{
				Addr:              params.Config.APIServer.HttpServerAddr,
				ReadHeaderTimeout: time.Duration(30) * time.Second,
			},
		}
		if params.Config.Tracing.OtelCollectorHTTP {
			server.server.Handler = otelhttp.NewHandler(e, "")
		} else {
			server.server.Handler = e
		}

		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				conn, err := grpc.NewClient(
					params.Config.APIServer.GrpcServerAddr,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
					grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
				)
				if err != nil {
					return err
				}

				gateway, err := grpcpkg.NewGateway(ctx, conn, params.Logger,
					grpcpkg.WithServerMuxOptions(
						runtime.WithErrorHandler(interceptors.HttpErrorHandler(params.Logger)),
						runtime.WithMetadata(interceptors.MetadataRequestPath()),
					),
					grpcpkg.WithHandlers(params.Register.HttpHandlers...),
				)
				if err != nil {
					return err
				}

				if params.Config.Tracing.OtelCollectorHTTP {
					server.echo.Any("/api/v1/*", echo.WrapHandler(httppkg.NewTraceparentWrapper(gateway)))
				} else {
					server.echo.Any("/api/v1/*", echo.WrapHandler(gateway))
				}

				return nil
			},
		})

		return server, nil
	}
}

func RunGatewayServer() func(logger *logger.Logger, server *GatewayServer) error {
	return func(logger *logger.Logger, server *GatewayServer) error {
		logger.Info("starting http server...")

		listener, err := net.Listen("tcp", server.ListenAddr)
		if err != nil {
			return fmt.Errorf("failed to listen %s: %v", server.ListenAddr, err)
		}

		go func() {
			err = server.server.Serve(listener)
			if err != nil && err != http.ErrServerClosed {
				logger.Fatal(err.Error())
			}
		}()

		logger.Info("http server listening...", zap.String("addr", server.ListenAddr))

		return nil
	}
}
