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
	ErrLessonAttemtNotFound       = errors.New("lesson attempt not found")
	ErrQuestionPageAttemtNotFound = errors.New("question page attempt not found")
	ErrAnswerNotFound             = errors.New("page answer not found")
)

func (lp *LpService) TryLesson(ctx context.Context, lesson *lpmodels.TryLesson) (*lpmodels.TryLessonResp, error) {
	const op = "internal.services.lp.attempts.TryLesson"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", lesson.UserID),
		slog.Int64("lesson_id", lesson.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "TryLesson")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", lesson.UserID),
		attribute.Int64("lesson_id", lesson.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(lesson); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("start trying lesson")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    lesson.UserID,
		ChannelID: lesson.ChannelID,
		PlanID:    lesson.PlanID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", lesson.UserID))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start trying
	log.Info("try lesson")
	span.AddEvent("started_trying_lesson")
	resp, err := lp.AttemptProvider.TryLesson(ctx, lesson)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrQuestionPageAttemtNotFound):
			log.Error("question page attempt not found", slog.Any("lesson_id", lesson.LessonID))
			return nil, fmt.Errorf("%s: %w", op, ErrQuestionPageAttemtNotFound)
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_trying_lesson")

	log.Info("trying successfully")

	return resp, nil
}

func (lp *LpService) UpdatePageAttempt(ctx context.Context, attempt *lpmodels.UpdatePageAttempt) (*lpmodels.UpdatePageAttemptResp, error) {
	const op = "internal.services.lp.attempts.UpdatePageAttempt"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", attempt.UserID),
		slog.Int64("lesson_attempt_id", attempt.LessonAttemptID),
	)

	_, span := tracer.LPtracer.Start(ctx, "UpdatePageAttempt")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", attempt.UserID),
		attribute.Int64("lesson_attempt_id", attempt.LessonAttemptID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(attempt); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("start updating page attempt")

	// Start check permissions
	span.AddEvent("checking_attempt_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckLessonAttemptPermissions(ctx, &lpmodels.LessonAttemptPermissions{
		UserID:          attempt.UserID,
		LessonAttemptID: attempt.LessonAttemptID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", attempt.UserID))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_attempt_permissons_for_user")

	// Start updating
	log.Info("updating question page attempt")
	span.AddEvent("started_update_question_page_attempt")
	resp, err := lp.AttemptProvider.UpdatePageAttempt(ctx, attempt)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrAnswerNotFound):
			log.Error("question page attempt not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrAnswerNotFound)
		default:
			log.Error("failed to update question page attempt", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_updating_question_page_attempt")

	log.Info("question page attempt updated successfully")

	return resp, nil
}

func (lp *LpService) CompleteLesson(ctx context.Context, lesson *lpmodels.CompleteLesson) (*lpmodels.CompleteLessonResp, error) {
	const op = "internal.services.lp.attempts.CompleteLesson"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", lesson.UserID),
		slog.Int64("lesson_attempt_id", lesson.LessonAttemptID),
	)

	_, span := tracer.LPtracer.Start(ctx, "CompleteLesson")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", lesson.UserID),
		attribute.Int64("lesson_attempt_id", lesson.LessonAttemptID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(lesson); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("start completing lesson attempt")

	// Start check permissions
	span.AddEvent("checking_attempt_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckLessonAttemptPermissions(ctx, &lpmodels.LessonAttemptPermissions{
		UserID:          lesson.UserID,
		LessonAttemptID: lesson.LessonAttemptID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", lesson.UserID))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_attempt_permissons_for_user")

	// Start complete
	log.Info("completing lesson attempt")
	span.AddEvent("started_complete_lesson_attempt")
	resp, err := lp.AttemptProvider.CompleteLesson(ctx, lesson)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrQuestionPageAttemtNotFound):
			log.Error("question page attempt not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrQuestionPageAttemtNotFound)
		default:
			log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_lesson_attempt")

	log.Info("lesson completed successfully")

	return resp, nil
}

func (lp *LpService) GetLessonAttempts(ctx context.Context, inputParams *lpmodels.GetLessonAttempts) (*lpmodels.GetLessonAttemptsResp, error) {
	const op = "internal.services.lp.attempts.GetLessonAttempts"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", inputParams.UserID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetLessonAttempts")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", inputParams.UserID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(inputParams); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("start getting lesson attempts")

	// Start Getting
	log.Info("getting lesson attempts")
	span.AddEvent("started_getting_lesson_attempts")
	resp, err := lp.AttemptProvider.GetLessonAttempts(ctx, inputParams)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrLessonAttemtNotFound):
			log.Error("lesson attempt not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrLessonAttemtNotFound)
		default:
			log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_lesson_attempts")

	log.Info("lesson attempts got successfully")

	return resp, nil
}
