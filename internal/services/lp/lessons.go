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
	ErrLessonNotFound = errors.New("lesson not found")
)

func (lp *LpService) CreateLesson(ctx context.Context, lesson *lpmodels.CreateLesson) (*lpmodels.CreateLessonResponse, error) {
	const op = "internal.services.lp.lessons.CreateLesson"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", lesson.CreatedBy),
		slog.String("new_lesson_name", lesson.Name),
	)

	_, span := tracer.LPtracer.Start(ctx, "CreateLesson")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", lesson.CreatedBy),
		attribute.String("new_lesson_name", lesson.Name),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(lesson); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("creating new lesson")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    lesson.CreatedBy,
		PlanID:    lesson.PlanID,
		ChannelID: lesson.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", lesson.CreatedBy))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start creating
	span.AddEvent("started_creating_lesson")
	resp, err := lp.LessonProvider.CreateLesson(ctx, lesson)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("name", lesson.Name))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to creating new lesson", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_creating_lesson")

	log.Info("lesson created successfully")

	return &lpmodels.CreateLessonResponse{
		ID:      resp.ID,
		Success: resp.Success,
	}, nil
}

func (lp *LpService) GetLesson(ctx context.Context, lesson *lpmodels.GetLesson) (*lpmodels.GetLessonResponse, error) {
	const op = "internal.services.lp.lessons.GetLesson"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", lesson.UserID),
		slog.Int64("lesson_id", lesson.LessonID),
		slog.Int64("plan_id", lesson.PlanID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetLesson")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", lesson.UserID),
		attribute.Int64("lesson_id", lesson.LessonID),
		attribute.Int64("plan_id", lesson.PlanID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(lesson); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.GetLessonResponse{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    lesson.UserID,
		ChannelID: lesson.ChannelID,
		PlanID:    lesson.PlanID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.GetLessonResponse{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", lesson.UserID))
		return &lpmodels.GetLessonResponse{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting lesson by id")
	span.AddEvent("started_getting_lesson_by_id")
	resp, err := lp.LessonProvider.GetLesson(ctx, lesson)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrLessonNotFound):
			log.Error("lesson not found", slog.Any("lesson_id", lesson.LessonID))
			return &lpmodels.GetLessonResponse{}, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		default:
			log.Error("failed to get lesson", slog.String("err", err.Error()))
			return &lpmodels.GetLessonResponse{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_lesson_by_id")

	log.Info("getting lesson successfully")

	return resp, nil
}

func (lp *LpService) GetLessons(ctx context.Context, inputParam *lpmodels.GetLessons) ([]lpmodels.GetLessonResponse, error) {
	const op = "internal.services.lp.lessons.GetLessons"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", inputParam.UserID),
		slog.Int64("plan_id", inputParam.PlanID),
		slog.Int64("channel_id", inputParam.ChannelID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetLessons")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", inputParam.UserID),
		attribute.Int64("plan_id", inputParam.PlanID),
		attribute.Int64("channel_id", inputParam.ChannelID),
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
	perm, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    inputParam.UserID,
		ChannelID: inputParam.ChannelID,
		PlanID:    inputParam.PlanID,
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
	log.Info("getting lessons")
	span.AddEvent("started_getting_lessons")
	resp, err := lp.LessonProvider.GetLessons(ctx, inputParam)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrLessonNotFound):
			log.Error("lessons not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to get lessons", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_lessons")

	log.Info("getting lessons successfully")

	return resp, nil
}

func (lp *LpService) UpdateLesson(ctx context.Context, updLesson *lpmodels.UpdateLesson) (*lpmodels.UpdateLessonResponse, error) {
	const op = "internal.services.lp.lessons.UpdateLesson"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", updLesson.LastModifiedBy),
		slog.Int64("lesson_id", updLesson.LessonID),
		slog.Int64("plan_id", updLesson.PlanID),
	)

	_, span := tracer.LPtracer.Start(ctx, "UpdateLesson")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", updLesson.LastModifiedBy),
		attribute.Int64("lesson_id", updLesson.LessonID),
		attribute.Int64("plan_id", updLesson.PlanID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(updLesson); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.UpdateLessonResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    updLesson.LastModifiedBy,
		PlanID:    updLesson.PlanID,
		ChannelID: updLesson.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.UpdateLessonResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", updLesson.LastModifiedBy))
		return &lpmodels.UpdateLessonResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start updating
	log.Info("updating lesson")
	span.AddEvent("started_update_lesson")
	resp, err := lp.LessonProvider.UpdateLesson(ctx, updLesson)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdateLessonResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrLessonNotFound):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdateLessonResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		default:
			log.Error("failed to update lesson", slog.String("err", err.Error()))
			return &lpmodels.UpdateLessonResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_updating_lesson")

	log.Info("lesson updated successfully")

	return resp, nil
}

func (lp *LpService) DeleteLesson(ctx context.Context, delLess *lpmodels.DeleteLesson) (*lpmodels.DeleteLessonResponse, error) {
	const op = "internal.services.lp.lessons.DeleteLesson"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", delLess.UserID),
		slog.Int64("lesson_id", delLess.LessonID),
		slog.Int64("plan_id", delLess.PlanID),
	)

	_, span := tracer.LPtracer.Start(ctx, "DeleteLesson")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", delLess.UserID),
		attribute.Int64("lesson_id", delLess.LessonID),
		attribute.Int64("plan_id", delLess.PlanID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(delLess); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.DeleteLessonResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    delLess.UserID,
		PlanID:    delLess.PlanID,
		ChannelID: delLess.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.DeleteLessonResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", delLess.UserID))
		return &lpmodels.DeleteLessonResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start deleting
	log.Info("deleting lesson")
	span.AddEvent("started_delete_lesson")
	resp, err := lp.LessonProvider.DeleteLesson(ctx, delLess)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrLessonNotFound):
			log.Error("lesson not found", slog.String("err", err.Error()))
			return &lpmodels.DeleteLessonResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		default:
			log.Error("failed to delete lesson", slog.String("err", err.Error()))
			return &lpmodels.DeleteLessonResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_deleting_lesson")

	log.Info("lesson deleted successfully")

	return resp, nil

}
