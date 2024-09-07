package interceptors

import (
	"context"
	"runtime/debug"

	"github.com/nekomeowww/factorio-rcon-api/pkg/apierrors"
	"github.com/nekomeowww/xo/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func PanicInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			r := recover()
			if r != nil {
				logger.Error("panicked", zap.Any("err", r), zap.Stack(string(debug.Stack())))
				err = apierrors.NewErrInternal().AsStatus()
				resp = nil
			}
		}()

		return handler(ctx, req)
	}
}
