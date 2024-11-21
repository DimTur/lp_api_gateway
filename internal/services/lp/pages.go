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
	ErrPageNotFound = errors.New("page not found")
)

func (lp *LpService) CreateImagePage(ctx context.Context, page *lpmodels.CreateImagePage) (*lpmodels.CreatePageResponse, error) {
	const op = "internal.services.lp.pages.CreateImagePage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", page.CreatedBy),
		slog.String("page_name", page.ImageName),
		slog.String("page_type", page.CreateBasePage.ContentType),
	)

	_, span := tracer.LPtracer.Start(ctx, "CreateImagePage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", page.CreatedBy),
		attribute.String("page_name", page.ImageName),
		attribute.String("page_type", page.CreateBasePage.ContentType),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(page); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("creating new image page")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    page.CreatedBy,
		PlanID:    page.PlanID,
		ChannelID: page.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", page.CreatedBy))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start creating
	span.AddEvent("started_creating_image_page")
	resp, err := lp.PageProvider.CreateImagePage(ctx, page)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("page_name", page.ImageName))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to creating new image page", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_creating_image_page")

	log.Info("image page created successfully")

	return &lpmodels.CreatePageResponse{
		ID:      resp.ID,
		Success: resp.Success,
	}, nil
}

func (lp *LpService) CreateVideoPage(ctx context.Context, page *lpmodels.CreateVideoPage) (*lpmodels.CreatePageResponse, error) {
	const op = "internal.services.lp.pages.CreateVideoPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", page.CreatedBy),
		slog.String("page_name", page.VideoName),
		slog.String("page_type", page.CreateBasePage.ContentType),
	)

	_, span := tracer.LPtracer.Start(ctx, "CreateVideoPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", page.CreatedBy),
		attribute.String("page_name", page.VideoName),
		attribute.String("page_type", page.CreateBasePage.ContentType),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(page); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("creating new video page")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    page.CreatedBy,
		PlanID:    page.PlanID,
		ChannelID: page.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", page.CreatedBy))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start creating
	span.AddEvent("started_creating_video_page")
	resp, err := lp.PageProvider.CreateVideoPage(ctx, page)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("page_name", page.VideoName))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to creating new video page", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_creating_video_page")

	log.Info("video page created successfully")

	return &lpmodels.CreatePageResponse{
		ID:      resp.ID,
		Success: resp.Success,
	}, nil
}

func (lp *LpService) CreatePdfPage(ctx context.Context, page *lpmodels.CreatePDFPage) (*lpmodels.CreatePageResponse, error) {
	const op = "internal.services.lp.pages.CreatePdfPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", page.CreatedBy),
		slog.String("page_name", page.PdfName),
		slog.String("page_type", page.CreateBasePage.ContentType),
	)

	_, span := tracer.LPtracer.Start(ctx, "CreatePdfPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", page.CreatedBy),
		attribute.String("page_name", page.PdfName),
		attribute.String("page_type", page.CreateBasePage.ContentType),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(page); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	log.Info("creating new pdf page")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    page.CreatedBy,
		PlanID:    page.PlanID,
		ChannelID: page.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", page.CreatedBy))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start creating
	span.AddEvent("started_creating_pdf_page")
	resp, err := lp.PageProvider.CreatePDFPage(ctx, page)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("invalid credentinals", slog.Any("page_name", page.PdfName))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to creating new pdf page", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_creating_pdf_page")

	log.Info("pdf page created successfully")

	return &lpmodels.CreatePageResponse{
		ID:      resp.ID,
		Success: resp.Success,
	}, nil
}

func (lp *LpService) GetImagePage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.ImagePage, error) {
	const op = "internal.services.lp.pages.GetImagePage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", page.UserID),
		slog.Int64("page_id", page.PageID),
		slog.Int64("lesson_id", page.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetImagePage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", page.UserID),
		attribute.Int64("page_id", page.PageID),
		attribute.Int64("lesson_id", page.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(page); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.ImagePage{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    page.UserID,
		ChannelID: page.ChannelID,
		PlanID:    page.PlanID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.ImagePage{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", page.UserID))
		return &lpmodels.ImagePage{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting image by id")
	span.AddEvent("started_getting_image_page_by_id")
	resp, err := lp.PageProvider.GetImagePage(ctx, page)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPageNotFound):
			log.Error("image page not found", slog.Any("page_id", page.LessonID))
			return &lpmodels.ImagePage{}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to get image page", slog.String("err", err.Error()))
			return &lpmodels.ImagePage{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_image page_by_id")

	log.Info("got image page successfully")

	return resp, nil
}

func (lp *LpService) GetVideoPage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.VideoPage, error) {
	const op = "internal.services.lp.pages.GetVideoPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", page.UserID),
		slog.Int64("page_id", page.PageID),
		slog.Int64("lesson_id", page.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetVideoPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", page.UserID),
		attribute.Int64("page_id", page.PageID),
		attribute.Int64("lesson_id", page.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(page); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.VideoPage{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    page.UserID,
		ChannelID: page.ChannelID,
		PlanID:    page.PlanID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.VideoPage{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", page.UserID))
		return &lpmodels.VideoPage{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting video page by id")
	span.AddEvent("started_getting_video_page_by_id")
	resp, err := lp.PageProvider.GetVideoPage(ctx, page)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPageNotFound):
			log.Error("video page not found", slog.Any("page_id", page.LessonID))
			return &lpmodels.VideoPage{}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to get video page", slog.String("err", err.Error()))
			return &lpmodels.VideoPage{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_video page_by_id")

	log.Info("got video page successfully")

	return resp, nil
}

func (lp *LpService) GetPDFPage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.PDFPage, error) {
	const op = "internal.services.lp.pages.GetPDFPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", page.UserID),
		slog.Int64("page_id", page.PageID),
		slog.Int64("lesson_id", page.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetPDFPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", page.UserID),
		attribute.Int64("page_id", page.PageID),
		attribute.Int64("lesson_id", page.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(page); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.PDFPage{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	p, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    page.UserID,
		ChannelID: page.ChannelID,
		PlanID:    page.PlanID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.PDFPage{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", page.UserID))
		return &lpmodels.PDFPage{}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting pdf page by id")
	span.AddEvent("started_getting_pdf_page_by_id")
	resp, err := lp.PageProvider.GetPDFPage(ctx, page)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPageNotFound):
			log.Error("pdf page not found", slog.Any("page_id", page.LessonID))
			return &lpmodels.PDFPage{}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to get pdf page", slog.String("err", err.Error()))
			return &lpmodels.PDFPage{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_pdf_page_by_id")

	log.Info("got pdf page successfully")

	return resp, nil
}

func (lp *LpService) GetPages(ctx context.Context, inputParams *lpmodels.GetPages) ([]lpmodels.BasePage, error) {
	const op = "internal.services.lp.pages.GetPages"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", inputParams.UserID),
		slog.Int64("lesson_id", inputParams.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "GetPages")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", inputParams.UserID),
		attribute.Int64("lesson_id", inputParams.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(inputParams); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	span.AddEvent("checking_permissons_for_user")
	perm, err := lp.PermissionsProvider.CheckCreaterOrLearnerAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    inputParams.UserID,
		ChannelID: inputParams.ChannelID,
		PlanID:    inputParams.PlanID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !perm {
		log.Info("permissions denied", slog.String("user_id", inputParams.UserID))
		return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	span.AddEvent("completed_checking_permissons_for_user")

	// Start getting
	log.Info("getting pages")
	span.AddEvent("started_getting_pages")
	resp, err := lp.PageProvider.GetPages(ctx, inputParams)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPageNotFound):
			log.Error("pages not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to get pages", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_getting_pages")

	log.Info("got pages successfully")

	return resp, nil
}

func (lp *LpService) UpdateImagePage(ctx context.Context, updIPage *lpmodels.UpdateImagePage) (*lpmodels.UpdatePageResponse, error) {
	const op = "internal.services.lp.pages.UpdateImagePage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", updIPage.LastModifiedBy),
		slog.Int64("page_id", updIPage.PlanID),
		slog.Int64("lesson_id", updIPage.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "UpdateImagePage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", updIPage.LastModifiedBy),
		attribute.Int64("page_id", updIPage.PlanID),
		attribute.Int64("lesson_id", updIPage.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(updIPage); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    updIPage.LastModifiedBy,
		PlanID:    updIPage.PlanID,
		ChannelID: updIPage.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", updIPage.LastModifiedBy))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start updating
	log.Info("updating image page")
	span.AddEvent("started_update_image_page")
	resp, err := lp.PageProvider.UpdateImagePage(ctx, updIPage)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrPageNotFound):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to update lesson", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_updating_image_page")

	log.Info("image page updated successfully")

	return resp, nil
}

func (lp *LpService) UpdateVideoPage(ctx context.Context, updIPage *lpmodels.UpdateVideoPage) (*lpmodels.UpdatePageResponse, error) {
	const op = "internal.services.lp.pages.UpdateVideoPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", updIPage.LastModifiedBy),
		slog.Int64("page_id", updIPage.PlanID),
		slog.Int64("lesson_id", updIPage.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "UpdateVideoPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", updIPage.LastModifiedBy),
		attribute.Int64("page_id", updIPage.PlanID),
		attribute.Int64("lesson_id", updIPage.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(updIPage); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    updIPage.LastModifiedBy,
		PlanID:    updIPage.PlanID,
		ChannelID: updIPage.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", updIPage.LastModifiedBy))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start updating
	log.Info("updating video page")
	span.AddEvent("started_update_video_page")
	resp, err := lp.PageProvider.UpdateVideoPage(ctx, updIPage)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrPageNotFound):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to update lesson", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_updating_video_page")

	log.Info("video page updated successfully")

	return resp, nil
}

func (lp *LpService) UpdatePDFPage(ctx context.Context, updIPage *lpmodels.UpdatePDFPage) (*lpmodels.UpdatePageResponse, error) {
	const op = "internal.services.lp.pages.UpdatePDFPage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", updIPage.LastModifiedBy),
		slog.Int64("page_id", updIPage.PlanID),
		slog.Int64("lesson_id", updIPage.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "UpdatePDFPage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", updIPage.LastModifiedBy),
		attribute.Int64("page_id", updIPage.PlanID),
		attribute.Int64("lesson_id", updIPage.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(updIPage); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    updIPage.LastModifiedBy,
		PlanID:    updIPage.PlanID,
		ChannelID: updIPage.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", updIPage.LastModifiedBy))
		return &lpmodels.UpdatePageResponse{
			ID:      0,
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start updating
	log.Info("updating pdf page")
	span.AddEvent("started_update_pdf_page")
	resp, err := lp.PageProvider.UpdatePDFPage(ctx, updIPage)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrInvalidCredentials):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case errors.Is(err, lpgrpc.ErrPageNotFound):
			log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to update lesson", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_updating_pdf_page")

	log.Info("pdf page updated successfully")

	return resp, nil
}

func (lp *LpService) DeletePage(ctx context.Context, delPage *lpmodels.DeletePage) (*lpmodels.DeletePageResponse, error) {
	const op = "internal.services.lp.pages.DeletePage"

	log := lp.Log.With(
		slog.String("op", op),
		slog.String("user_id", delPage.UserID),
		slog.Int64("page_id", delPage.PageID),
		slog.Int64("lesson_id", delPage.LessonID),
	)

	_, span := tracer.LPtracer.Start(ctx, "DeletePage")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", delPage.UserID),
		attribute.Int64("page_id", delPage.PageID),
		attribute.Int64("lesson_id", delPage.LessonID),
	)

	// Validation
	span.AddEvent("validation_started")
	if err := lp.Validator.Struct(delPage); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return &lpmodels.DeletePageResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	span.AddEvent("validation_completed")

	// Start check permissions
	p, err := lp.PermissionsProvider.CheckCreatorOrAdminAndSharePermissions(ctx, &permissions.CheckPerm{
		UserID:    delPage.UserID,
		PlanID:    delPage.PlanID,
		ChannelID: delPage.ChannelID,
	})
	if err != nil {
		log.Error("can't check permissions", slog.String("err", err.Error()))
		return &lpmodels.DeletePageResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}
	if !p {
		log.Info("permissions denied", slog.String("user_id", delPage.UserID))
		return &lpmodels.DeletePageResponse{
			Success: false,
		}, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
	}

	// Start deleting
	log.Info("deleting page")
	span.AddEvent("started_delete_page")
	resp, err := lp.PageProvider.DeletePage(ctx, delPage)
	if err != nil {
		switch {
		case errors.Is(err, lpgrpc.ErrPageNotFound):
			log.Error("page not found", slog.String("err", err.Error()))
			return &lpmodels.DeletePageResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			log.Error("failed to delete page", slog.String("err", err.Error()))
			return &lpmodels.DeletePageResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}
	span.AddEvent("completed_deleting_page")

	log.Info("lesson deleted successfully")

	return resp, nil
}
