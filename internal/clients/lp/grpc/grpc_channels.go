package lpgrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrChannelNotFound    = errors.New("channel not found")

	ErrInternal = errors.New("internal error")
)

func (c *Client) CreateChannel(ctx context.Context, newChannel *lpmodels.CreateChannel) (*lpmodels.CreateChannelResponse, error) {
	const op = "lp.grpc.CreateChannel"

	channel := &lpv1.CreateChannelRequest{
		Name:           newChannel.Name,
		Description:    newChannel.Description,
		CreatedBy:      newChannel.CreatedBy,
		LastModifiedBy: newChannel.CreatedBy,
	}

	resp, err := c.api.CreateChannel(ctx, channel)
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("invalid arguments", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.CreateChannelResponse{
		ID:      resp.Id,
		Success: true,
	}, nil
}

func (c *Client) GetChannel(ctx context.Context, channel *lpmodels.GetChannel) (*lpmodels.GetChannelResponse, error) {
	const op = "lp.grpc.GetChannel"

	resp, err := c.api.GetChannel(ctx, &lpv1.GetChannelRequest{
		Id: channel.ChannelID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("channel not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	channelResponse := &lpmodels.GetChannelResponse{
		Id:             resp.Channel.Id,
		Name:           resp.Channel.Name,
		Description:    resp.Channel.Description,
		CreatedBy:      resp.Channel.CreatedBy,
		LastModifiedBy: resp.Channel.LastModifiedBy,
		CreatedAt:      resp.Channel.CreatedAt.AsTime(),
		Modified:       resp.Channel.Modified.AsTime(),
		Plans:          make([]*lpmodels.Plan, len(resp.Channel.Plans)),
	}

	for i, plan := range resp.Channel.Plans {
		channelResponse.Plans[i] = &lpmodels.Plan{
			Id:             plan.Id,
			Name:           plan.Name,
			Description:    plan.Description,
			CreatedBy:      plan.CreatedBy,
			LastModifiedBy: plan.LastModifiedBy,
			IsPublished:    plan.IsPublished,
			Public:         plan.Public,
			CreatedAt:      plan.CreatedAt.AsTime(),
			Modified:       plan.Modified.AsTime(),
		}
	}

	return channelResponse, nil
}
