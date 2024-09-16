package lpgrpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api lpv1.LearningPlatformClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "lp.grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// TODO: secure conn
	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: lpv1.NewLearningPlatformClient(cc),
		log: log,
	}, nil
}

func (c *Client) CreateChannel(ctx context.Context,
	name string,
	description string,
	userID int64,
	public bool) (*lpv1.CreateChannelResponse, error) {
	const op = "lp.grpc.CreateChannel"

	channel := &lpv1.Channel{
		Name:        name,
		Description: description,
		CreatedBy:   userID,
		Public:      public,
	}

	resp, err := c.api.CreateChannel(ctx, &lpv1.CreateChannelRequest{
		Channel: channel,
	})
	if err != nil {
		c.log.Error("received error from auth grpc service", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (c *Client) GetChannel(ctx context.Context, channelID int64) (*lpv1.GetChannelResponse, error) {
	const op = "lp.grpc.GetChannel"

	resp, err := c.api.GetChannel(ctx, &lpv1.GetChannelRequest{
		ChannelId: channelID,
	})
	if err != nil {
		c.log.Error("received error from auth grpc service", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

// InterceptorLogger adapts slog logger to iterceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
