package libs

import (
	"context"
	"io"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/nekomeowww/factorio-rcon-api/v2/internal/configs"
	"github.com/nekomeowww/xo/logger"
)

type NewOtelParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *configs.Config
	Logger    *logger.Logger
}

type Otel struct {
	Trace  *trace.TracerProvider
	Metric *metric.MeterProvider
}

func NewOtel() func(params NewOtelParams) (*Otel, error) {
	return func(params NewOtelParams) (*Otel, error) {
		prop := propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)

		otel.SetTextMapPropagator(prop)

		o := &Otel{}

		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				{
					var err error
					var spanExporter trace.SpanExporter

					if params.Config.Tracing.OtelCollectorHTTP {
						params.Logger.Info("configured to use otlp collector for tracing", zap.String("protocol", "http/protobuf"), zap.String("endpoint", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")))

						spanExporter, err = otlptrace.New(ctx, otlptracehttp.NewClient())
						if err != nil {
							return err
						}
					} else if params.Config.Tracing.OtelStdoutEnabled {
						params.Logger.Info("configured to use stdout exporter for tracing")

						spanExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
						if err != nil {
							return err
						}
					} else {
						params.Logger.Info("configured to disable stdout exporter for tracing")

						spanExporter, err = stdouttrace.New(stdouttrace.WithWriter(io.Discard))
						if err != nil {
							return err
						}
					}

					tracerProvider := trace.NewTracerProvider(trace.WithBatcher(spanExporter))
					otel.SetTracerProvider(tracerProvider)

					o.Trace = tracerProvider

					params.Lifecycle.Append(fx.Hook{
						OnStop: func(ctx context.Context) error {
							return tracerProvider.Shutdown(ctx)
						},
					})
				}

				{
					var err error
					var metricExporter metric.Exporter

					if params.Config.Tracing.OtelCollectorHTTP {
						params.Logger.Info("configured to use otlp collector for metric", zap.String("protocol", "http/protobuf"), zap.String("endpoint", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")))

						metricExporter, err = otlpmetrichttp.New(ctx)
						if err != nil {
							return err
						}
					} else if params.Config.Tracing.OtelStdoutEnabled {
						params.Logger.Info("configured to use stdout exporter for metric")

						metricExporter, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
						if err != nil {
							return err
						}
					} else {
						params.Logger.Info("configured to disable stdout exporter for metric")

						metricExporter, err = stdoutmetric.New(stdoutmetric.WithWriter(io.Discard))
						if err != nil {
							return err
						}
					}

					meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(metricExporter)))
					otel.SetMeterProvider(meterProvider)

					o.Metric = meterProvider

					params.Lifecycle.Append(fx.Hook{
						OnStop: func(ctx context.Context) error {
							return meterProvider.Shutdown(ctx)
						},
					})
				}

				err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
				if err != nil {
					return err
				}

				return nil
			},
		})

		return o, nil
	}
}
