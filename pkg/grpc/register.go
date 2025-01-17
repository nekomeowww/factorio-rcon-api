package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type HTTPHandler = func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error
type GRPCServiceRegister func(s reflection.GRPCServer)

type Register struct {
	HTTPHandlers []HTTPHandler
	GrpcServices []GRPCServiceRegister
	EchoHandlers map[string]map[string]echo.HandlerFunc
}

func NewRegister() *Register {
	return &Register{
		HTTPHandlers: make([]HTTPHandler, 0),
		GrpcServices: make([]GRPCServiceRegister, 0),
		EchoHandlers: make(map[string]map[string]echo.HandlerFunc),
	}
}

func (r *Register) RegisterHTTPHandler(handler HTTPHandler) {
	r.HTTPHandlers = append(r.HTTPHandlers, handler)
}

func (r *Register) RegisterHTTPHandlers(handlers []HTTPHandler) {
	r.HTTPHandlers = append(r.HTTPHandlers, handlers...)
}

func (r *Register) RegisterGrpcService(serviceRegister GRPCServiceRegister) {
	r.GrpcServices = append(r.GrpcServices, serviceRegister)
}

func (r *Register) RegisterEchoHandler(path string, method string, handler echo.HandlerFunc) {
	if _, ok := r.EchoHandlers[path]; !ok {
		r.EchoHandlers[path] = make(map[string]echo.HandlerFunc)
	}

	r.EchoHandlers[path][method] = handler
}
