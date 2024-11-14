package permissions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPermissionDenied   = errors.New("permissions denied")
	ErrInternal           = errors.New("internal error")
)

func (p *PermissionsService) IsGroupAdmin(ctx context.Context, uIsGroupAdmin *ssomodels.IsGroupAdmin) (bool, error) {
	const op = "internal.services.permissions.permissions.IsGroupAdmin"

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", uIsGroupAdmin.UserID),
		slog.String("learning_group_id", uIsGroupAdmin.LgID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "IsGroupAdmin")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := p.validator.Struct(uIsGroupAdmin); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("user_id", uIsGroupAdmin.UserID))
	span.SetAttributes(attribute.String("learning_group_id", uIsGroupAdmin.LgID))

	log.Info("start checking permissions for user")
	span.AddEvent("checking_permissions_for_user")
	// checks that user is group admin
	log.Info("checks that user is group admin")
	span.AddEvent("checking_user_is_group_admin")
	isGroupAdmin, err := p.lgPermissionsProvider.IsGroupAdmin(ctx, &ssomodels.IsGroupAdmin{
		UserID: uIsGroupAdmin.UserID,
		LgID:   uIsGroupAdmin.LgID,
	})
	if err != nil {
		log.Error("can't check that user is group admin", slog.String("err", err.Error()))
	}

	// If user creator - returns true immediately
	if isGroupAdmin.IsGroupAdmin {
		span.AddEvent("completed_checking_permissions_for_user")
		span.SetAttributes(attribute.String("user_id", uIsGroupAdmin.UserID))
		span.SetAttributes(attribute.String("learning_group_id", uIsGroupAdmin.LgID))
		return true, nil
	} else {
		span.AddEvent("completed_checking_permissions_for_user")
		span.SetAttributes(attribute.String("user_id", uIsGroupAdmin.UserID))
		span.SetAttributes(attribute.String("learning_group_id", uIsGroupAdmin.LgID))
		return false, nil
	}
}

func (p *PermissionsService) CheckCreaterOrLearnerAndSharePermissions(ctx context.Context, perm *CheckPerm) (bool, error) {
	const op = "internal.services.permissions.permissions.CheckCreaterOrLearnerAndSharePermissions"

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", perm.UserID),
		slog.Int64("channel_id", perm.ChannelID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "CheckCreaterOrLearnerAndSharePermissions")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := p.validator.Struct(perm); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("user_id", perm.UserID))
	span.SetAttributes(attribute.Int64("channel_id", perm.ChannelID))

	log.Info("start checking permissions for user")
	span.AddEvent("checking_permissions_for_user")
	// checks that channel created by user
	log.Info("checks that channel created by user")
	span.AddEvent("checking_channel_created_by_user")
	isChannelCreatorResp, err := p.channelPermissionsProvider.IsChannelCreator(ctx, &lpmodels.IsChannelCreator{
		UserID:    perm.UserID,
		ChannelID: perm.ChannelID,
	})
	if err != nil {
		log.Error("can't check that user is channel creator", slog.String("err", err.Error()))
	}

	// If user creator - returns true immediately
	if isChannelCreatorResp.IsCreator {
		span.AddEvent("completed_checking_permissions_for_user")
		span.SetAttributes(attribute.String("user_id", perm.UserID))
		span.SetAttributes(attribute.Int64("channel_id", perm.ChannelID))
		return true, nil
	}

	// Get learning groups ids relevant for user, where he is learner
	log.Info("getting learning groups ids where user is learner")
	span.AddEvent("get_learning_groups_ids_to_check_permissions")
	lgUserIsLearner, err := p.lgPermissionsProvider.UserIsLearnerIn(ctx, &ssomodels.UserIsLearnerIn{
		UserID: perm.UserID,
	})
	if err != nil {
		log.Error("can't get learning group ids where user is learner", slog.String("err", err.Error()))
	}

	// Save to Redis
	if err := p.redisPermissionsProvider.SaveLgUser(ctx, perm.UserID, lgUserIsLearner); err != nil {
		log.Error("can't save to redis learning group ids where user is learner", slog.String("err", err.Error()))
	}

	// Get learning groups ids with which the channel has been sharing
	log.Info("getting learning groups ids with which the channel has been sharing")
	span.AddEvent("get_learning_groups_ids_to_check_permissions")
	lgShareWithChannel, err := p.channelPermissionsProvider.LerningGroupsShareWithChannel(ctx, &lpmodels.LerningGroupsShareWithChannel{
		ChannelID: perm.ChannelID,
	})
	if err != nil {
		log.Error("can't get learning group ids with which the channel has been sharing", slog.String("err", err.Error()))
	}

	// Save to Redis
	if err := p.redisPermissionsProvider.SaveLgShareWithChannel(ctx, perm.ChannelID, lgShareWithChannel); err != nil {
		log.Error("can't save learning group ids with which the channel has been sharing", slog.String("err", err.Error()))
	}

	hasIntersection, err := p.redisPermissionsProvider.CheckGroupsIntersection(ctx, perm.UserID, perm.ChannelID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	return hasIntersection, nil
}

func (p *PermissionsService) CheckCreatorOrAdminAndSharePermissions(ctx context.Context, perm *CheckPerm) (bool, error) {
	const op = "internal.services.permissions.permissions.CheckCreatorOrAdminAndSharePermissions"

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", perm.UserID),
		slog.Int64("channel_id", perm.ChannelID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "CheckCreatorOrAdminAndSharePermissions")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := p.validator.Struct(perm); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("user_id", perm.UserID))
	span.SetAttributes(attribute.Int64("channel_id", perm.ChannelID))

	log.Info("start checking permissions for user")
	span.AddEvent("checking_permissions_for_user")
	// checks that channel created by user
	log.Info("checks that channel created by user")
	span.AddEvent("checking_channel_created_by_user")
	isChannelCreatorResp, err := p.channelPermissionsProvider.IsChannelCreator(ctx, &lpmodels.IsChannelCreator{
		UserID:    perm.UserID,
		ChannelID: perm.ChannelID,
	})
	if err != nil {
		log.Error("can't check that user is channel creator", slog.String("err", err.Error()))
	}

	// If user creator - returns true immediately
	if isChannelCreatorResp.IsCreator {
		span.AddEvent("completed_checking_permissions_for_user")
		span.SetAttributes(attribute.String("user_id", perm.UserID))
		span.SetAttributes(attribute.Int64("channel_id", perm.ChannelID))
		return true, nil
	}

	// Get learning groups ids relevant for user, where he is admin
	log.Info("getting learning groups ids where user is admin")
	span.AddEvent("get_learning_groups_ids_to_check_permissions")
	lgUserIsAdmin, err := p.lgPermissionsProvider.UserIsGroupAdminIn(ctx, &ssomodels.UserIsGroupAdminIn{
		UserID: perm.UserID,
	})
	if err != nil {
		log.Error("can't get learning group ids where user is admin", slog.String("err", err.Error()))
	}

	// Save to Redis
	if err := p.redisPermissionsProvider.SaveLgUser(ctx, perm.UserID, lgUserIsAdmin); err != nil {
		log.Error("can't save to redis learning group ids where user is admin", slog.String("err", err.Error()))
	}

	// Get learning groups ids with which the channel has been sharing
	log.Info("getting learning groups ids with which the channel has been sharing")
	span.AddEvent("get_learning_groups_ids_to_check_permissions")
	lgShareWithChannel, err := p.channelPermissionsProvider.LerningGroupsShareWithChannel(ctx, &lpmodels.LerningGroupsShareWithChannel{
		ChannelID: perm.ChannelID,
	})
	if err != nil {
		log.Error("can't get learning group ids with which the channel has been sharing", slog.String("err", err.Error()))
	}

	// Save to Redis
	if err := p.redisPermissionsProvider.SaveLgShareWithChannel(ctx, perm.ChannelID, lgShareWithChannel); err != nil {
		log.Error("can't save learning group ids with which the channel has been sharing", slog.String("err", err.Error()))
	}

	hasIntersection, err := p.redisPermissionsProvider.CheckGroupsIntersection(ctx, perm.UserID, perm.ChannelID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	return hasIntersection, nil
}
