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
		ChannelId:        channel.ChannelID,
		LearningGroupIds: channel.LearningGroupIds,
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
		CreatedAt:      resp.Channel.CreatedAt,
		Modified:       resp.Channel.Modified,
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
			CreatedAt:      plan.CreatedAt,
			Modified:       plan.Modified,
		}
	}

	return channelResponse, nil
}

func (c *Client) GetChannels(ctx context.Context, inputParam *lpmodels.GetChannels) ([]lpmodels.Channel, error) {
	const op = "lp.grpc.GetChannels"

	resp, err := c.api.GetChannels(ctx, &lpv1.GetChannelsRequest{
		LearningGroupIds: inputParam.LearningGroupIds,
		Limit:            inputParam.Limit,
		Offset:           inputParam.Offset,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("channels not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	var chanResp []lpmodels.Channel
	for _, channel := range resp.Channels {
		chanResp = append(chanResp, lpmodels.Channel{
			ID:             channel.Id,
			Name:           channel.Name,
			Description:    channel.Description,
			CreatedBy:      channel.CreatedBy,
			LastModifiedBy: channel.LastModifiedBy,
			CreatedAt:      channel.CreatedAt,
			Modified:       channel.Modified,
		})
	}

	return chanResp, nil
}

func (c *Client) UpdateChannel(ctx context.Context, updChannel *lpmodels.UpdateChannel) (*lpmodels.UpdateChannelResponse, error) {
	const op = "lp.grpc.UpdateChannel"

	resp, err := c.api.UpdateChannel(ctx, &lpv1.UpdateChannelRequest{
		UserId:       updChannel.UserID,
		AdminInLgIds: updChannel.AdminInLgIds,
		ChannelId:    updChannel.ChannelID,
		Name:         updChannel.Name,
		Description:  updChannel.Description,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.UpdateChannelResponse{
		ID:      resp.Id,
		Success: true,
	}, nil
}

func (c *Client) DeleteChannel(ctx context.Context, delChannel *lpmodels.DelChByID) (*lpmodels.DelChByIDResp, error) {
	const op = "lp.grpc.DeleteChannel"

	resp, err := c.api.DeleteChannel(ctx, &lpv1.DeleteChannelRequest{
		ChannelId:    delChannel.ChannelID,
		AdminInLgIds: delChannel.AdminInLgIds,
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

	return &lpmodels.DelChByIDResp{
		Success: resp.Success,
	}, nil
}

func (c *Client) ShareChannelToGroup(ctx context.Context, s *lpmodels.SharingChannel) (*lpmodels.SharingChannelResp, error) {
	const op = "lp.grpc.ShareChannelToGroup"

	resp, err := c.api.ShareChannelToGroup(ctx, &lpv1.ShareChannelToGroupRequest{
		ChannelId:  s.ChannelID,
		LgroupsIds: s.LGroupIDs,
		CreatedBy:  s.CreatedBy,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.SharingChannelResp{
		Success: resp.Success,
	}, nil
}
