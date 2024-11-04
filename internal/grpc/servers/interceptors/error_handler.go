package interceptors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nekomeowww/factorio-rcon-api/v2/apis/jsonapi"
	"github.com/nekomeowww/factorio-rcon-api/v2/pkg/apierrors"
	"github.com/nekomeowww/xo/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleStatusError(s *status.Status, logger *logger.Logger, err error) *apierrors.ErrResponse {
	switch s.Code() { //nolint
	case codes.InvalidArgument:
		if len(s.Details()) > 0 {
			break
		}

		return apierrors.NewErrInvalidArgument().WithDetail(s.Message()).AsResponse()
	case codes.Unimplemented:
		logger.Error("unimplemented error", zap.Error(err))
		return apierrors.NewErrNotFound().WithDetail("route not found or method not allowed").AsResponse()
	case codes.Internal:
		var errorCaller *jsonapi.ErrorCaller

		if len(s.Details()) > 1 {
			errorCaller, _ = s.Details()[1].(*jsonapi.ErrorCaller)
		}

		fields := []zap.Field{zap.Error(err)}
		if errorCaller != nil {
			fields = append(fields, zap.String("file", fmt.Sprintf("%s:%d", errorCaller.File, errorCaller.Line)))
			fields = append(fields, zap.String("function", errorCaller.Function))
		}

		logger.Error("internal error", fields...)

		return apierrors.NewErrInternal().AsResponse()
	case codes.NotFound:
		if len(s.Details()) > 0 {
			break
		}

		logger.Error("unimplemented error", zap.Error(err))

		return apierrors.NewErrNotFound().WithDetail("route not found or method not allowed").AsResponse()
	case codes.DeadlineExceeded:
		if len(s.Details()) > 0 {
			break
		}

		return apierrors.NewErrTimeout().AsResponse()
	default:
		break
	}

	errResp := apierrors.NewErrResponse()
	if len(s.Details()) > 0 {
		detail, ok := s.Details()[0].(*jsonapi.ErrorObject)
		if ok {
			errResp = errResp.WithError(&apierrors.Error{
				ErrorObject: detail,
			})
		}
	}

	return errResp
}

func handleError(logger *logger.Logger, err error) *apierrors.ErrResponse {
	if s, ok := status.FromError(err); ok {
		return handleStatusError(s, logger, err)
	}

	logger.Error("unknown error (probably unhandled)", zap.Error(err))

	return apierrors.NewErrInternal().AsResponse()
}

func HttpErrorHandler(logger *logger.Logger) func(ctx context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, writer http.ResponseWriter, _ *http.Request, err error) {
	return func(ctx context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, writer http.ResponseWriter, _ *http.Request, err error) {
		if err != nil {
			errResp := handleError(logger, err)

			b, _ := json.Marshal(errResp)

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(errResp.HttpStatus())

			_, _ = writer.Write(b)
		}
	}
}
