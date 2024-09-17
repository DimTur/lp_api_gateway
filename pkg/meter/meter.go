package meter

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	metr "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	ReqMeter = otel.Meter("requests-meter")

	AllReqCount, _           = ReqMeter.Int64Counter("requests_total", metr.WithDescription("Total number of requests"))
	SignUpReqCount, _        = ReqMeter.Int64Counter("requests_sign_up", metr.WithDescription("Sign up number of requests"))
	SignInReqCount, _        = ReqMeter.Int64Counter("requests_sign_in", metr.WithDescription("Sign in number of requests"))
	CreateChannelReqCount, _ = ReqMeter.Int64Counter("requests_create_channel", metr.WithDescription("Create Channel number of requests"))
	GetChannelReqCount, _    = ReqMeter.Int64Counter("requests_get_channel", metr.WithDescription("Get Channel in number of requests"))
)

func InitMeter(ctx context.Context, serviceName string) (*metric.MeterProvider, error) {
	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String("service.name", serviceName),
	))
	if err != nil {
		return nil, err
	}

	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}
