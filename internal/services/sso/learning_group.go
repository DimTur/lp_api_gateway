package ssoservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	ssogrpc "github.com/DimTur/lp_api_gateway/internal/clients/sso/grpc"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidGroupID     = errors.New("invalid group id")
	ErrGroupExists        = errors.New("group already exists")
	ErrGroupNotFound      = errors.New("group not found")
	ErrPermissionDenied   = errors.New("you don't have permissions")

	ErrInternal = errors.New("internal error")
)

func (sso *SsoService) CreateLearningGroup(ctx context.Context, newLg *ssomodels.CreateLearningGroup) (*ssomodels.CreateLearningGroupResp, error) {
	const op = "internal.services.sso.learning_group.CreateLearningGroup"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("new_learning_group_name", newLg.Name),
	)

	_, span := tracer.AuthTracer.Start(ctx, "CreateLearningGroup")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(newLg); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("new_learning_group_name", newLg.Name))

	log.Info("creating new learning group")

	// Start creating
	span.AddEvent("started_creating_learning_group")
	resp, err := sso.LgProvider.CreateLearningGroup(ctx, newLg)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("name", newLg.Name))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, ssogrpc.ErrGroupExists):
			log.Error("learning group already exists", slog.Any("name", newLg.Name))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupExists)
		default:
			log.Error("failed to creating new learning group", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_creating_learning_group")
	span.SetAttributes(attribute.String("name", newLg.Name))

	log.Info("learning group created successfully")

	return &ssomodels.CreateLearningGroupResp{
		Success: resp.Success,
	}, nil
}

func (sso *SsoService) GetLearningGroupByID(ctx context.Context, lgID *ssomodels.GetLgByID) (*ssomodels.GetLgByIDResp, error) {
	const op = "internal.services.sso.learning_group.GetLearningGroupByID"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_id", lgID.UserID),
		slog.String("learning_group_id", lgID.LgId),
	)

	_, span := tracer.AuthTracer.Start(ctx, "GetLearningGroupByID")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(lgID); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("learning_group_id", lgID.LgId))

	log.Info("getting learning group by id")

	// Start getting
	span.AddEvent("started_getting_learning_group_by_id")
	resp, err := sso.LgProvider.GetLearningGroupByID(ctx, lgID)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrPermissionDenied):
			log.Error("permissions denied", slog.Any("learning_group_id", lgID.LgId))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, ssogrpc.ErrGroupNotFound):
			log.Error("learning group not found", slog.Any("learning_group_id", lgID.LgId))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		default:
			log.Error("failed to get learning group", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_learning_group_by_id")
	span.SetAttributes(attribute.String("learning_group_id", lgID.LgId))

	log.Info("getting learning group successfully")

	return resp, nil
}

func (sso *SsoService) UpdateLearningGroup(ctx context.Context, updFields *ssomodels.UpdateLearningGroup) (*ssomodels.UpdateLearningGroupResp, error) {
	const op = "internal.services.sso.learning_group.UpdateLearningGroup"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_id", updFields.UserID),
		slog.String("learning_group_id", updFields.LgId),
	)

	_, span := tracer.AuthTracer.Start(ctx, "UpdateLearningGroup")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(updFields); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("learning_group_id", updFields.LgId))

	log.Info("updating learning group by id")

	// Start updating
	span.AddEvent("updating_learning_group")
	resp, err := sso.LgProvider.UpdateLearningGroup(ctx, updFields)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrPermissionDenied):
			log.Error("permissions denied", slog.Any("learning_group_id", updFields.LgId))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, ssogrpc.ErrGroupNotFound):
			log.Error("learning group not found", slog.Any("learning_group_id", updFields.LgId))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("learning_group_id", updFields.LgId))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to update learning group", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_update_learning_group")
	span.SetAttributes(attribute.String("learning_group_id", updFields.LgId))

	log.Info("learning group updated successfully")

	return &ssomodels.UpdateLearningGroupResp{
		Success: resp.Success,
	}, nil
}

func (sso *SsoService) DeleteLearningGroup(ctx context.Context, lgID *ssomodels.DelLgByID) (*ssomodels.DelLgByIDResp, error) {
	const op = "internal.services.sso.learning_group.DeleteLearningGroup"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_id", lgID.UserID),
		slog.String("learning_group_id", lgID.LgID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "DeleteLearningGroup")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(lgID); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("learning_group_id", lgID.LgID))

	log.Info("deleting learning group by id")

	// Start deleting
	span.AddEvent("deleting_learning_group")
	resp, err := sso.LgProvider.DeleteLearningGroup(ctx, lgID)
	log.Info("learning group deleted successfully")
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrPermissionDenied):
			log.Error("permissions denied", slog.Any("learning_group_id", lgID.LgID))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, ssogrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("learning_group_id", lgID.LgID))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to delete learning group", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_delete_learning_group")
	span.SetAttributes(attribute.String("learning_group_id", lgID.LgID))

	return &ssomodels.DelLgByIDResp{
		Success: resp.Success,
	}, nil
}

func (sso *SsoService) GetLearningGroups(ctx context.Context, uID *ssomodels.GetLGroups) (*ssomodels.GetLGroupsResp, error) {
	const op = "internal.services.sso.learning_group.GetLearningGroups"

	log := sso.Log.With(
		slog.String("op", op),
		slog.String("user_id", uID.UserID),
	)

	_, span := tracer.AuthTracer.Start(ctx, "GetLearningGroups")
	defer span.End()

	// Validation
	span.AddEvent("validation_started")
	if err := sso.Validator.Struct(uID); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")
	span.SetAttributes(attribute.String("user_id", uID.UserID))

	log.Info("getting learning groups")

	// Start getting
	span.AddEvent("started_getting_learning_groups")
	resp, err := sso.LgProvider.GetLearningGroups(ctx, uID)
	if err != nil {
		switch {
		case errors.Is(err, ssogrpc.ErrPermissionDenied):
			log.Error("permissions denied", slog.Any("user_id", uID.UserID))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, ssogrpc.ErrGroupNotFound):
			log.Error("learning group not found", slog.Any("user_id", uID.UserID))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		default:
			log.Error("failed to get learning group", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_learning_groups")
	span.SetAttributes(attribute.String("user_id", uID.UserID))

	log.Info("got learning groups successfully")

	return resp, nil
}
