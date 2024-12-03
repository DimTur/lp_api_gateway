package lpservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	lpgrpc "github.com/DimTur/lp_api_gateway/internal/clients/lp/grpc"
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/services/permissions"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrPlanNotFound = errors.New("plan not found")
)

func (lp *LpService) CreatePlan(ctx context.Context, plan *lpmodels.CreatePlan) (*lpmodels.CreatePlanResponse, error) {
	const op = "internal.services.lp.plans.CreatePlan"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", plan.CreatedBy),
		slog.String("new_plan_name", plan.Name),
	)

	_, span := tracer.LPtracer.Start(ctx, "CreatePlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", plan.CreatedBy),
		attribute.String("new_plan_name", plan.Name),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(plan); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("creating new plan")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    plan.CreatedBy,
		ChannelID: plan.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", plan.CreatedBy))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start creating
	span.AddEvent("started_creating_plan")
	resp, err := lp.PlanProvider.CreatePlan(ctx, plan)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("name", plan.Name))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to creating new plan", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_creating_plan")

	log.Info("plan created successfully")

	return &lpmodels.CreatePlanResponse{
		ID:      resp.ID,
		Success: resp.Success,
	}, nil
}

func (lp *LpService) GetPlan(ctx context.Context, plan *lpmodels.GetPlan) (*lpmodels.GetPlanResponse, error) {
	const op = "internal.services.lp.plans.GetPlan"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", plan.UserID),
		slog.Int64("plan_id", plan.PlanID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetPlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", plan.UserID),
		attribute.Int64("plan_id", plan.PlanID),
		attribute.Int64("channel_id", plan.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(plan); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.GetPlanResponse{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    plan.UserID,
		PlanID:    plan.PlanID,
		ChannelID: plan.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.GetPlanResponse{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", plan.UserID))
		return &lpmodels.GetPlanResponse{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting plan by id")
	span.AddEvent("started_getting_plan_by_id")
	resp, err := lp.PlanProvider.GetPlan(ctx, plan)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPlanNotFound):
			log.Error("plan not found", slog.Any("plan_id", plan.PlanID))
			return &lpmodels.GetPlanResponse{}, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		default:
			log.Error("failed to get plan", slog.String("err", err.Error()))
			return &lpmodels.GetPlanResponse{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_plan_by_id")

	log.Info("getting plan successfully")

	return resp, nil
}

func (lp *LpService) GetPlans(ctx context.Context, inputParam *lpmodels.GetPlans) ([]lpmodels.GetPlanResponse, error) {
	const op = "internal.services.lp.plans.GetPlans"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", inputParam.UserID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetPlans")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", inputParam.UserID),
		attribute.Int64("channel", inputParam.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(inputParam); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check admin permissions
	span.AddEvent("checking_permissons_for_user")
	adminPerm, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    inputParam.UserID,
		ChannelID: inputParam.ChannelID,
	})
	fmt.Println("err", err)
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		// return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	fmt.Println("adminPerm", adminPerm)
	if adminPerm {
		log.Info("getting plans")
		span.AddEvent("started_getting_plans")
		resp, err := lp.PlanProvider.GetPlansForGroupAdmin(ctx, inputParam)
		if err != nil {
			switch {
			case errors.Is(err, lpgrpc.ErrPlanNotFound):
				log.Info("plans not found", slog.String("err", err.Error()))
				return nil, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
			case errors.Is(err, lpgrpc.ErrInvalidCredentials):
				log.Info("bad request", slog.String("err", err.Error()))
				return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
			default:
				log.Info("failed to get plans", slog.String("err", err.Error()))
				return nil, fmt.Errorf("%s: %w", op, ErrInternal)
			}
		}
		span.AddEvent("completed_getting_plans")

		return resp, nil
	}

	// Start check learner permissions
	perm, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    inputParam.UserID,
		ChannelID: inputParam.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !perm {
		log.Info("permissions denied", slog.String("user_id", inputParam.UserID))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting plans")
	span.AddEvent("started_getting_plans")
	resp, err := lp.PlanProvider.GetPlans(ctx, inputParam)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPlanNotFound):
			log.Info("plans not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Info("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Info("failed to get plans", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_plans")

	log.Info("getting plans successfully")

	return resp, nil
}

func (lp *LpService) GetPlansForGroupAdmin(ctx context.Context, inputParam *lpmodels.GetPlans) ([]lpmodels.GetPlanResponse, error) {
	const op = "internal.services.lp.plans.GetPlansAll"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", inputParam.UserID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetPlansAll")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", inputParam.UserID),
		attribute.Int64("channel", inputParam.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(inputParam); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	perm, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    inputParam.UserID,
		ChannelID: inputParam.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !perm {
		log.Info("permissions denied", slog.String("user_id", inputParam.UserID))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting plans")
	span.AddEvent("started_getting_plans")
	resp, err := lp.PlanProvider.GetPlans(ctx, inputParam)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPlanNotFound):
			log.Error("plans not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to get plans", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_plans")

	log.Info("getting plans successfully")

	return resp, nil
}

func (lp *LpService) UpdatePlan(ctx context.Context, updPlan *lpmodels.UpdatePlan) (*lpmodels.UpdatePlanResponse, error) {
	const op = "internal.services.lp.plans.UpdatePlan"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", updPlan.LastModifiedBy),
		slog.Int64("plan_id", updPlan.PlanID),
	)

	_, span := tracer.LPtracer.Start(ctx, "UpdatePlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", updPlan.LastModifiedBy),
		attribute.Int64("plan_id", updPlan.PlanID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(updPlan); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.UpdatePlanResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    updPlan.LastModifiedBy,
		PlanID:    updPlan.PlanID,
		ChannelID: updPlan.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.UpdatePlanResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", updPlan.LastModifiedBy))
		return &lpmodels.UpdatePlanResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start updating
	log.Info("updating plan")
	span.AddEvent("started_update_plan")
	resp, err := lp.PlanProvider.UpdatePlan(ctx, updPlan)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePlanResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrPlanNotFound):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePlanResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		default:
			log.Error("failed to update plan", slog.String("err", err.Error()))
			return &lpmodels.UpdatePlanResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_updating_plan")

	log.Info("channel plan successfully")

	return resp, nil
}

func (lp *LpService) DeletePlan(ctx context.Context, delPlan *lpmodels.DelPlan) (*lpmodels.DelPlanResponse, error) {
	const op = "internal.services.lp.plans.DeletePlan"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", delPlan.UserID),
		slog.Int64("plan_id", delPlan.PlanID),
		slog.Int64("channel_id", delPlan.ChannelID),
	)

	_, span := tracer.LPtracer.Start(ctx, "DeletePlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", delPlan.UserID),
		attribute.Int64("plan_id", delPlan.PlanID),
		attribute.Int64("channel_id", delPlan.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(delPlan); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.DelPlanResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    delPlan.UserID,
		PlanID:    delPlan.PlanID,
		ChannelID: delPlan.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.DelPlanResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", delPlan.UserID))
		return &lpmodels.DelPlanResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start deleting
	log.Info("deleting plan")
	span.AddEvent("started_delete_plan")
	resp, err := lp.PlanProvider.DeletePlan(ctx, delPlan)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPlanNotFound):
			log.Error("plan not found", slog.String("err", err.Error()))
			return &lpmodels.DelPlanResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		default:
			log.Error("failed to delete plan", slog.String("err", err.Error()))
			return &lpmodels.DelPlanResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_deleting_plan")

	log.Info("plan deleted successfully")

	return resp, nil
}

func (lp *LpService) SharePlanWithUser(ctx context.Context, sharePlanWithUser *lpmodels.SharePlan) (*lpmodels.SharingPlanResp, error) {
	const op = "internal.services.lp.plans.SharePlanWithUser"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", sharePlanWithUser.UserID),
		slog.Int64("plan_id", sharePlanWithUser.PlanID),
		slog.Int64("channel_id", sharePlanWithUser.ChannelID),
	)

	_, span := tracer.LPtracer.Start(ctx, "ShareChannelToGroup")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", sharePlanWithUser.UserID),
		attribute.Int64("plan_id", sharePlanWithUser.PlanID),
		attribute.Int64("channel_id", sharePlanWithUser.ChannelID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(sharePlanWithUser); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.SharingPlanResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    sharePlanWithUser.UserID,
		PlanID:    sharePlanWithUser.PlanID,
		ChannelID: sharePlanWithUser.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.SharingPlanResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", sharePlanWithUser.UserID))
		return &lpmodels.SharingPlanResp{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start sharing
	log.Info("sharing plan")
	span.AddEvent("started_share_plan")
	resp, err := lp.PlanProvider.SharePlanWithUser(ctx, sharePlanWithUser)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.SharingPlanResp{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to share plan", slog.String("err", err.Error()))
			return &lpmodels.SharingPlanResp{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_deleting_plan")

	log.Info("plan shared successfully")

	return resp, nil
}
