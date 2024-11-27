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
	"go.opentelemetry.io/otel/trace"
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
		slog.Int64("plan_id", perm.PlanID),
		slog.Int64("channel_id", perm.ChannelID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "CheckCreaterOrLearnerAndSharePermissions")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", perm.UserID),
		attribute.Int64("plan_id", perm.PlanID),
		attribute.Int64("channel_id", perm.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := p.validator.Struct(perm); err != nil {
		span.AddEvent("validation_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// checks that channel created by user
	log.Info("checks that channel created by user")
	span.AddEvent("check_channel_creator_started")
	isChannelCreatorResp, err := p.channelPermissionsProvider.IsChannelCreator(ctx, &lpmodels.IsChannelCreator{
		UserID:    perm.UserID,
		ChannelID: perm.ChannelID,
	})
	if err != nil {
		span.AddEvent("check_channel_creator_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't check that user is channel creator", slog.String("err", err.Error()))
	}

	// If user creator - returns true immediately
	if isChannelCreatorResp.IsCreator {
		span.AddEvent("user_is_channel_creator", trace.WithAttributes(
			attribute.String("user_id", perm.UserID),
			attribute.Int64("channel_id", perm.ChannelID),
		))
		return true, nil
	}
	span.AddEvent("check_channel_creator_completed")

	// Get learning groups ids relevant for user, where he is learner
	log.Info("getting learning groups ids where user is learner")
	span.AddEvent("fetch_learning_groups_for_user_started")
	lgUserIsLearner, err := p.lgPermissionsProvider.UserIsLearnerIn(ctx, &ssomodels.UserIsLearnerIn{
		UserID: perm.UserID,
	})
	if err != nil {
		span.AddEvent("fetch_learning_groups_for_user_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't get learning group ids where user is learner", slog.String("err", err.Error()))
	}
	span.AddEvent("fetch_learning_groups_for_user_completed")

	// Save to Redis
	span.AddEvent("save_learning_groups_to_redis_started")
	if err := p.redisPermissionsProvider.SaveLgUser(ctx, perm.UserID, lgUserIsLearner); err != nil {
		span.AddEvent("save_learning_groups_to_redis_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't save to redis learning group ids where user is learner", slog.String("err", err.Error()))
	}
	span.AddEvent("save_learning_groups_to_redis_completed")

	// Get learning groups ids with which the channel has been sharing
	log.Info("getting learning groups ids with which the channel has been sharing")
	span.AddEvent("fetch_shared_learning_groups_started")
	lgShareWithChannel, err := p.channelPermissionsProvider.LerningGroupsShareWithChannel(ctx, &lpmodels.LerningGroupsShareWithChannel{
		ChannelID: perm.ChannelID,
	})
	if err != nil {
		span.AddEvent("fetch_shared_learning_groups_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't get learning group ids with which the channel has been sharing", slog.String("err", err.Error()))
	}
	span.AddEvent("fetch_shared_learning_groups_completed")

	// Save to Redis
	span.AddEvent("save_shared_groups_to_redis_started")
	if err := p.redisPermissionsProvider.SaveLgShareWithChannel(ctx, perm.ChannelID, lgShareWithChannel); err != nil {
		span.AddEvent("save_shared_groups_to_redis_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't save learning group ids with which the channel has been sharing", slog.String("err", err.Error()))
	}
	span.AddEvent("save_shared_groups_to_redis_completed")

	// Check intersection
	span.AddEvent("check_groups_intersection_started")
	hasLgIntersection, err := p.redisPermissionsProvider.CheckGroupsIntersection(ctx, perm.UserID, perm.ChannelID)
	if err != nil {
		span.AddEvent("check_groups_intersection_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("check_groups_intersection_completed", trace.WithAttributes(attribute.Bool("has_intersection", hasLgIntersection)))

	if !hasLgIntersection {
		span.AddEvent("permissions_denied", trace.WithAttributes(attribute.String("reason", "no_group_intersection")))
		log.Warn("permissions denied for", slog.String("user_id", perm.UserID))
		return false, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	if perm.PlanID == 0 {
		span.AddEvent("permissions_granted_by_groups_intersection")
		return hasLgIntersection, nil
	}

	// Checking user id with which the plan has been sharing
	span.AddEvent("check_plan_permissions_started")
	isShare, err := p.planPermissionsProvider.IsUserShareWithPlan(ctx, &IsUserShareWithPlan{
		UserID: perm.UserID,
		PlanID: perm.PlanID,
	})
	if err != nil {
		span.AddEvent("check_plan_permissions_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't check user id with which the plan has been sharing", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("check_plan_permissions_completed", trace.WithAttributes(attribute.Bool("is_shared", isShare.IsShare)))

	if isShare.IsShare {
		span.AddEvent("permissions_granted_by_plan_sharing")
		return true, nil
	}

	span.AddEvent("permissions_denied", trace.WithAttributes(attribute.String("reason", "no_access_to_plan")))
	log.Warn("permissions denied for", slog.String("user_id", perm.UserID))
	return false, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
}

func (p *PermissionsService) CheckCreatorOrAdminAndSharePermissions(ctx context.Context, perm *CheckPerm) (bool, error) {
	const op = "internal.services.permissions.permissions.CheckCreatorOrAdminAndSharePermissions"

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", perm.UserID),
		slog.Int64("plan_id", perm.PlanID),
		slog.Int64("channel_id", perm.ChannelID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "CheckCreatorOrAdminAndSharePermissions")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", perm.UserID),
		attribute.Int64("plan_id", perm.PlanID),
		attribute.Int64("channel_id", perm.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := p.validator.Struct(perm); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// checks that channel created by user
	log.Info("checks that channel created by user")
	span.AddEvent("checking_channel_created_by_user")
	isChannelCreatorResp, err := p.channelPermissionsProvider.IsChannelCreator(ctx, &lpmodels.IsChannelCreator{
		UserID:    perm.UserID,
		ChannelID: perm.ChannelID,
	})
	if err != nil {
		span.AddEvent("check_channel_creator_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't check that user is channel creator", slog.String("err", err.Error()))
	}

	// If user creator - returns true immediately
	if isChannelCreatorResp.IsCreator {
		span.AddEvent("user_is_channel_creator", trace.WithAttributes(
			attribute.String("user_id", perm.UserID),
			attribute.Int64("channel_id", perm.ChannelID),
		))
		return true, nil
	}
	span.AddEvent("check_channel_creator_completed")

	// Get learning groups ids relevant for user, where he is admin
	log.Info("getting learning groups ids where user is admin")
	span.AddEvent("fetch_learning_groups_for_user_started")
	lgUserIsAdmin, err := p.lgPermissionsProvider.UserIsGroupAdminIn(ctx, &ssomodels.UserIsGroupAdminIn{
		UserID: perm.UserID,
	})
	if err != nil {
		span.AddEvent("fetch_learning_groups_for_user_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't get learning group ids where user is admin", slog.String("err", err.Error()))
	}
	span.AddEvent("fetch_learning_groups_for_user_completed")

	// Save to Redis
	span.AddEvent("save_learning_groups_to_redis_started")
	if err := p.redisPermissionsProvider.SaveLgUser(ctx, perm.UserID, lgUserIsAdmin); err != nil {
		span.AddEvent("save_learning_groups_to_redis_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't save to redis learning group ids where user is admin", slog.String("err", err.Error()))
	}
	span.AddEvent("save_learning_groups_to_redis_completed")

	// Get learning groups ids with which the channel has been sharing
	log.Info("getting learning groups ids with which the channel has been sharing")
	span.AddEvent("fetch_shared_learning_groups_started")
	lgShareWithChannel, err := p.channelPermissionsProvider.LerningGroupsShareWithChannel(ctx, &lpmodels.LerningGroupsShareWithChannel{
		ChannelID: perm.ChannelID,
	})
	if err != nil {
		span.AddEvent("fetch_shared_learning_groups_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't get learning group ids with which the channel has been sharing", slog.String("err", err.Error()))
	}
	span.AddEvent("fetch_shared_learning_groups_completed")

	// Save to Redis
	span.AddEvent("save_shared_groups_to_redis_started")
	if err := p.redisPermissionsProvider.SaveLgShareWithChannel(ctx, perm.ChannelID, lgShareWithChannel); err != nil {
		span.AddEvent("save_shared_groups_to_redis_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't save learning group ids with which the channel has been sharing", slog.String("err", err.Error()))
	}
	span.AddEvent("save_shared_groups_to_redis_completed")

	// Check intersection
	span.AddEvent("check_groups_intersection_started")
	hasLgIntersection, err := p.redisPermissionsProvider.CheckGroupsIntersection(ctx, perm.UserID, perm.ChannelID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("check_groups_intersection_completed", trace.WithAttributes(attribute.Bool("has_intersection", hasLgIntersection)))

	if !hasLgIntersection {
		span.AddEvent("permissions_denied", trace.WithAttributes(attribute.String("reason", "no_group_intersection")))
		log.Warn("permissions denied for", slog.String("user_id", perm.UserID))
		return false, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	if perm.PlanID == 0 {
		span.AddEvent("permissions_granted_by_groups_intersection")
		return hasLgIntersection, nil
	}

	// Checking user id with which the plan has been sharing
	span.AddEvent("check_plan_permissions_started")
	isShare, err := p.planPermissionsProvider.IsUserShareWithPlan(ctx, &IsUserShareWithPlan{
		UserID: perm.UserID,
		PlanID: perm.PlanID,
	})
	if err != nil {
		span.AddEvent("check_plan_permissions_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't check user id with which the plan has been sharing", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("check_plan_permissions_completed", trace.WithAttributes(attribute.Bool("is_shared", isShare.IsShare)))

	if isShare.IsShare {
		span.AddEvent("permissions_granted_by_plan_sharing")
		return true, nil
	}

	span.AddEvent("permissions_denied", trace.WithAttributes(attribute.String("reason", "no_access_to_plan")))
	log.Warn("permissions denied for", slog.String("user_id", perm.UserID))
	return false, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
}

func (p *PermissionsService) CheckLessonAttemptPermissions(ctx context.Context, userAtt *lpmodels.LessonAttemptPermissions) (bool, error) {
	const op = "internal.services.permissions.permissions.CheckCreatorOrAdminAndSharePermissions"

	log := p.log.With(
		slog.String("op", op),
		slog.String("user_id", userAtt.UserID),
		slog.Int64("lesson_attempt_id", userAtt.LessonAttemptID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "CheckLessonAttemptPermissions")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", userAtt.UserID),
		attribute.Int64("plan_id", userAtt.LessonAttemptID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := p.validator.Struct(userAtt); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Checks that the user is the creator of the attempt
	log.Info("checks that attempt created by user")
	span.AddEvent("checking_attemp_created_by_user")
	perm, err := p.attemptPermissionsProvider.CheckLessonAttemptPermissions(ctx, &lpmodels.LessonAttemptPermissions{
		UserID:          userAtt.UserID,
		LessonAttemptID: userAtt.LessonAttemptID,
	})
	if err != nil {
		span.AddEvent("check_attempt_creator_failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Error("can't check that user is attempt creator", slog.String("err", err.Error()))
	}

	return perm, nil
}
