package lpservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	lpgrpc "github.com/DimTur/lp_api_gateway/internal/clients/lp/grpc"
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/internal/services/permissions"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidChannelID   = errors.New("invalid channel id")
	ErrChannelExitsts     = errors.New("channel already exists")
	ErrChannelNotFound    = errors.New("channel not found")

	ErrPermissionDenied = errors.New("permissions denied")
	ErrInternal         = errors.New("internal error")
)

func (lp *LpService) CreateChannel(ctx context.Context, newChannel *lpmodels.CreateChannel) (*lpmodels.CreateChannelResponse, error) {
	const op = "internal.services.lp.channels.CreateChannel"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", newChannel.CreatedBy),
		slog.String("new_channel_name", newChannel.Name),
	)

	_, span := tracer.LPtracer.Start(ctx, "CreateChannel")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(newChannel); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("user_id", newChannel.CreatedBy))
	span.SetAttributes(attribute.String("new_channel_name", newChannel.Name))

	log.Info("creating new channel")

	// Start check permissions
	p, err := lp.PermissionsProvider.IsGroupAdmin(ctx, &ssomodels.IsGroupAdmin{
		UserID: newChannel.CreatedBy,
		LgID:   newChannel.LearningGroupId,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", newChannel.CreatedBy))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start creating
	span.AddEvent("started_creating_channel")
	resp, err := lp.ChannelProvider.CreateChannel(ctx, newChannel)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("name", newChannel.Name))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to creating new channel", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_creating_channel")
	span.SetAttributes(attribute.String("user_id", newChannel.CreatedBy))
	span.SetAttributes(attribute.String("new_channel_name", newChannel.Name))

	log.Info("channel created successfully")

	return &lpmodels.CreateChannelResponse{
		ID:      resp.ID,
		Success: resp.Success,
	}, nil
}

func (lp *LpService) GetChannel(ctx context.Context, channel *lpmodels.GetChannel) (*lpmodels.GetChannelResponse, error) {
	const op = "internal.services.lp.channels.GetChannel"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", channel.UserID),
		slog.Int64("channel_id", channel.ChannelID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetChannel")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", channel.UserID),
		attribute.Int64("channel_id", channel.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(channel); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.GetChannelResponse{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    channel.UserID,
		ChannelID: channel.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.GetChannelResponse{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", channel.UserID))
		return &lpmodels.GetChannelResponse{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start getting
	log.Info("getting channel by id")
	span.AddEvent("started_getting_channel_by_id")
	resp, err := lp.ChannelProvider.GetChannel(ctx, channel)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrChannelNotFound):
			log.Error("channel not found", slog.Any("channel_id", channel.ChannelID))
			return &lpmodels.GetChannelResponse{}, fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		default:
			log.Error("failed to get channel", slog.String("err", err.Error()))
			return &lpmodels.GetChannelResponse{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_channel_by_id")

	log.Info("getting channel successfully")

	return resp, nil
}

func (lp *LpService) GetChannels(ctx context.Context, inputParam *lpmodels.GetChannels) ([]lpmodels.Channel, error) {
	const op = "internal.services.lp.channels.GetChannels"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", inputParam.UserID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetChannels")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", inputParam.UserID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(inputParam); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Get learning groups ids relevant for user
	log.Info("getting learning groups ids where user is learner")
	span.AddEvent("checking_permissons_for_user")
	perm, err := lp.LgServiceProvider.UserIsLearnerIn(ctx, &ssomodels.UserIsLearnerIn{
		UserID: inputParam.UserID,
	})
	if err != nil {
		log.Error("permissions denied", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting channels")
	span.AddEvent("started_getting_channels")
	resp, err := lp.ChannelProvider.GetChannels(ctx, &lpmodels.GetChannelsFull{
		UserID:           inputParam.UserID,
		LearningGroupIds: perm,
		Limit:            inputParam.Limit,
		Offset:           inputParam.Offset,
	})
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrChannelNotFound):
			log.Error("channels not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to get channel", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_channels")

	log.Info("getting channels successfully")

	return resp, nil
}

func (lp *LpService) UpdateChannel(ctx context.Context, updChannel *lpmodels.UpdateChannel) (*lpmodels.UpdateChannelResponse, error) {
	const op = "internal.services.lp.channels.UpdateChannel"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", updChannel.UserID),
		slog.Int64("channel_id", updChannel.ChannelID),
	)

	_, span := tracer.LPtracer.Start(ctx, "UpdateChannel")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", updChannel.UserID),
		attribute.Int64("channel_id", updChannel.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(updChannel); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.UpdateChannelResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    updChannel.UserID,
		ChannelID: updChannel.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.UpdateChannelResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", updChannel.UserID))
		return &lpmodels.UpdateChannelResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start updating
	log.Info("updating channel")
	span.AddEvent("started_update_channel")
	resp, err := lp.ChannelProvider.UpdateChannel(ctx, updChannel)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdateChannelResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrPlanNotFound):
			log.Error("plan not found", slog.String("err", err.Error()))
			return &lpmodels.UpdateChannelResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		default:
			log.Error("failed to update channel", slog.String("err", err.Error()))
			return &lpmodels.UpdateChannelResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_updating_channel")

	log.Info("channel updated successfully")

	return resp, nil
}

func (lp *LpService) DeleteChannel(ctx context.Context, delChannel *lpmodels.DelChByID) (*lpmodels.DelChByIDResp, error) {
	const op = "internal.services.lp.channels.DeleteChannel"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", delChannel.UserID),
		slog.Int64("channel_id", delChannel.ChannelID),
	)

	_, span := tracer.LPtracer.Start(ctx, "DeleteChannel")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(delChannel); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.DelChByIDResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    delChannel.UserID,
		ChannelID: delChannel.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.DelChByIDResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", delChannel.UserID))
		return &lpmodels.DelChByIDResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start deleting
	log.Info("deleting channel")
	span.AddEvent("started_delete_channel")
	resp, err := lp.ChannelProvider.DeleteChannel(ctx, delChannel)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrChannelNotFound):
			log.Error("channel not found", slog.String("err", err.Error()))
			return &lpmodels.DelChByIDResp{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		default:
			log.Error("failed to delete channel", slog.String("err", err.Error()))
			return &lpmodels.DelChByIDResp{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_deleting_channel")

	log.Info("channel deleted successfully")

	return resp, nil
}

func (lp *LpService) ShareChannelToGroup(ctx context.Context, s *lpmodels.SharingChannel) (*lpmodels.SharingChannelResp, error) {
	const op = "internal.services.lp.channels.ShareChannelToGroup"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", s.UserID),
		slog.Int64("channel_id", s.ChannelID),
	)

	_, span := tracer.LPtracer.Start(ctx, "ShareChannelToGroup")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(s); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.SharingChannelResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("user_id", s.UserID))
	span.SetAttributes(attribute.Int64("channel_id", s.ChannelID))

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    s.UserID,
		ChannelID: s.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.SharingChannelResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", s.UserID))
		return &lpmodels.SharingChannelResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start sharing
	log.Info("sharing channel")
	span.AddEvent("started_share_channel")
	resp, err := lp.ChannelProvider.ShareChannelToGroup(ctx, &lpmodels.SharingChannel{
		UserID:    s.UserID,
		ChannelID: s.ChannelID,
		LGroupIDs: s.LGroupIDs,
	})
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.SharingChannelResp{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to share channel", slog.String("err", err.Error()))
			return &lpmodels.SharingChannelResp{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_deleting_channel")
	span.SetAttributes(attribute.String("user_id", s.UserID))
	span.SetAttributes(attribute.Int64("channel_id", s.ChannelID))

	log.Info("channel shared successfully")

	return resp, nil
}
