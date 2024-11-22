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
	ErrQuestionNotFound = errors.New("question not found")
)

func (lp *LpService) CreateQuestionPage(ctx context.Context, question *lpmodels.CreateQuestionPage) (*lpmodels.CreatePageResponse, error) {
	const op = "internal.services.lp.questions.CreateQuestionPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", question.CreatedBy),
	)

	_, span := tracer.LPtracer.Start(ctx, "CreateQuestionPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", question.CreatedBy),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(question); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("creating new question page")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    question.CreatedBy,
		PlanID:    question.PlanID,
		ChannelID: question.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", question.CreatedBy))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start creating
	span.AddEvent("started_creating_question_page")
	resp, err := lp.QuestionProvider.CreateQuestionPage(ctx, question)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to creating new question page", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_creating_question_page")

	log.Info("question page created successfully")

	return &lpmodels.CreatePageResponse{
		ID:      resp.ID,
		Success: resp.Success,
	}, nil
}

func (lp *LpService) GetQuestionPage(ctx context.Context, question *lpmodels.GetPage) (*lpmodels.GetQuestionPage, error) {
	const op = "internal.services.lp.questions.GetQuestionPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", question.UserID),
		slog.Int64("page_id", question.PageID),
		slog.Int64("lesson_id", question.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetQuestionPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", question.UserID),
		attribute.Int64("page_id", question.PageID),
		attribute.Int64("lesson_id", question.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(question); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.GetQuestionPage{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    question.UserID,
		ChannelID: question.ChannelID,
		PlanID:    question.PlanID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.GetQuestionPage{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", question.UserID))
		return &lpmodels.GetQuestionPage{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting question by id")
	span.AddEvent("started_getting_image_page_by_id")
	resp, err := lp.QuestionProvider.GetQuestionPage(ctx, question)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrQuestionNotFound):
			log.Error("question page not found", slog.Any("page_id", question.PageID))
			return &lpmodels.GetQuestionPage{}, fmt.Errorf("%s: %w", op, ErrQuestionNotFound)
		default:
			log.Error("failed to get question page", slog.String("err", err.Error()))
			return &lpmodels.GetQuestionPage{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_question_page_by_id")

	log.Info("got question page successfully")

	return resp, nil
}

func (lp *LpService) UpdateQuestionPage(ctx context.Context, updQust *lpmodels.UpdateQuestionPage) (*lpmodels.UpdatePageResponse, error) {
	const op = "internal.services.lp.questions.UpdateQuestionPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", updQust.LastModifiedBy),
		slog.Int64("page_id", updQust.PlanID),
		slog.Int64("lesson_id", updQust.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "UpdateQuestionPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", updQust.LastModifiedBy),
		attribute.Int64("page_id", updQust.PlanID),
		attribute.Int64("lesson_id", updQust.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(updQust); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    updQust.LastModifiedBy,
		PlanID:    updQust.PlanID,
		ChannelID: updQust.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", updQust.LastModifiedBy))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start updating
	log.Info("updating question page")
	span.AddEvent("started_update_question_page")
	resp, err := lp.QuestionProvider.UpdateQuestionPage(ctx, updQust)
	fmt.Println("ERRORRRRRRRRRRRRRRR", err)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrQuestionNotFound):
			log.Error("question not found", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrQuestionNotFound)
		default:
			log.Error("failed to update question page", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_updating_question_page")

	log.Info("question page updated successfully")

	return resp, nil
}
