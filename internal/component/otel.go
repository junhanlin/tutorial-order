package component

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/fx"
	shared "tutorial.io/tutorial-order/internal"
)

func NewOtelClient(
	lc fx.Lifecycle,
	config *shared.Config,
) *otlptrace.Client {
	otelClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			exporter, err := otlptrace.New(context.Background(), otelClient)
			if err != nil {
				log.WithError(err).Error("failed to initialize exporter")
			}
			tp := tracesdk.NewTracerProvider(
				tracesdk.WithBatcher(exporter),
				tracesdk.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceName(config.TraceServiceName),
				)),
			)
			otel.SetTracerProvider(tp)
			otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return otelClient.Stop(ctx)
		},
	})
	return &otelClient
}
