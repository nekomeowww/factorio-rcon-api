package grpc

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nekomeowww/xo/logger"
	"google.golang.org/grpc"
)

type gatewayOptions struct {
	serverMuxOptions []runtime.ServeMuxOption
	handlers         []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error
}

type GatewayCallOption func(*gatewayOptions)

func WithServerMuxOptions(opts ...runtime.ServeMuxOption) GatewayCallOption {
	return func(o *gatewayOptions) {
		o.serverMuxOptions = append(o.serverMuxOptions, opts...)
	}
}

func WithHandlers(handlers ...func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error) GatewayCallOption {
	return func(o *gatewayOptions) {
		o.handlers = append(o.handlers, handlers...)
	}
}

func NewGateway(
	ctx context.Context,
	conn *grpc.ClientConn,
	logger *logger.Logger,
	callOpts ...GatewayCallOption,
) (http.Handler, error) {
	opts := &gatewayOptions{}

	for _, f := range callOpts {
		f(opts)
	}

	mux := runtime.NewServeMux(opts.serverMuxOptions...)

	for _, f := range opts.handlers {
		err := f(ctx, mux, conn)
		if err != nil {
			return nil, err
		}
	}

	return mux, nil
}
