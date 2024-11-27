package pageshandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/handlers/utils"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	lpservice "github.com/DimTur/lp_api_gateway/internal/services/lp"
	"github.com/DimTur/lp_api_gateway/pkg/meter"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPageNotFound       = errors.New("page not found")
)

type LPService interface {
	CreateImagePage(ctx context.Context, page *lpmodels.CreateImagePage) (*lpmodels.CreatePageResponse, error)
	CreateVideoPage(ctx context.Context, page *lpmodels.CreateVideoPage) (*lpmodels.CreatePageResponse, error)
	CreatePdfPage(ctx context.Context, page *lpmodels.CreatePDFPage) (*lpmodels.CreatePageResponse, error)
	GetImagePage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.ImagePage, error)
	GetVideoPage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.VideoPage, error)
	GetPDFPage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.PDFPage, error)
	GetPages(ctx context.Context, inputParams *lpmodels.GetPages) ([]lpmodels.BasePage, error)
	UpdateImagePage(ctx context.Context, updIPage *lpmodels.UpdateImagePage) (*lpmodels.UpdatePageResponse, error)
	UpdateVideoPage(ctx context.Context, updIPage *lpmodels.UpdateVideoPage) (*lpmodels.UpdatePageResponse, error)
	UpdatePDFPage(ctx context.Context, updIPage *lpmodels.UpdatePDFPage) (*lpmodels.UpdatePageResponse, error)
	DeletePage(ctx context.Context, delPage *lpmodels.DeletePage) (*lpmodels.DeletePageResponse, error)
}

// CreateImagePage godoc
// @Summary      Create a new image page
// @Description  This endpoint allows users to create a new image page with the specified data.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        pageshandler.CreateImagePageRequest body pageshandler.CreateImagePageRequest true "Image page creation parameters"
// @Success      201 {object} pageshandler.CreatePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/image_page [post]
// @Security ApiKeyAuth
func CreateImagePage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.CreateImagePage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateImagePageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		req, err := utils.DecodeRequestBody[CreateImagePageRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.CreateImagePage(r.Context(), &lpmodels.CreateImagePage{
			CreateBasePage: lpmodels.CreateBasePage{
				LessonID:  lessonID,
				PlanID:    planID,
				ChannelID: channelID,
				CreatedBy: uID,
			},
			ImageFileUrl: req.ImageFileUrl,
			ImageName:    req.ImageName,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("invalid credentials", slog.Any("image_page", req.ImageName))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to create image page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to create image page"))
				return
			}
		}

		log.Info("image page created", slog.Int64("id", resp.ID))

		render.JSON(w, r, CreatePageResponse{
			Response: response.OK(),
			PageID:   resp.ID,
		})
		w.WriteHeader(http.StatusCreated)
	}
}

// CreateVideoPage godoc
// @Summary      Create a new video page
// @Description  This endpoint allows users to create a new video page with the specified data.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        pageshandler.CreateVideoPageRequest body pageshandler.CreateVideoPageRequest true "Video page creation parameters"
// @Success      201 {object} pageshandler.CreatePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/video_page [post]
// @Security ApiKeyAuth
func CreateVideoPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.CreateVideoPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateVideoPageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		req, err := utils.DecodeRequestBody[CreateVideoPageRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.CreateVideoPage(r.Context(), &lpmodels.CreateVideoPage{
			CreateBasePage: lpmodels.CreateBasePage{
				LessonID:  lessonID,
				PlanID:    planID,
				ChannelID: channelID,
				CreatedBy: uID,
			},
			VideoFileUrl: req.VideoFileUrl,
			VideoName:    req.VideoName,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("invalid credentials", slog.Any("video_page", req.VideoName))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to create video page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to create video page"))
				return
			}
		}

		log.Info("video page created", slog.Int64("id", resp.ID))

		render.JSON(w, r, CreatePageResponse{
			Response: response.OK(),
			PageID:   resp.ID,
		})
		w.WriteHeader(http.StatusCreated)
	}
}

// CreatePDFPage godoc
// @Summary      Create a new pdf page
// @Description  This endpoint allows users to create a new pdf page with the specified data.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        pageshandler.CreatePDFPageRequest body pageshandler.CreatePDFPageRequest true "PDF page creation parameters"
// @Success      201 {object} pageshandler.CreatePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pdf_page [post]
// @Security ApiKeyAuth
func CreatePDFPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.CreatePDFPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreatePDFPageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		req, err := utils.DecodeRequestBody[CreatePDFPageRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.CreatePdfPage(r.Context(), &lpmodels.CreatePDFPage{
			CreateBasePage: lpmodels.CreateBasePage{
				LessonID:  lessonID,
				PlanID:    planID,
				ChannelID: channelID,
				CreatedBy: uID,
			},
			PdfFileUrl: req.PdfFileUrl,
			PdfName:    req.PdfName,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("invalid credentials", slog.Any("pdf_page", req.PdfName))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to create pdf page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to create pdf page"))
				return
			}
		}

		log.Info("pdf page created", slog.Int64("id", resp.ID))

		render.JSON(w, r, CreatePageResponse{
			Response: response.OK(),
			PageID:   resp.ID,
		})
		w.WriteHeader(http.StatusCreated)
	}
}

// GetImagePage godoc
// @Summary      Get image page information
// @Description  This endpoint returns image information by ID.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Success      200 {object} pageshandler.GetImagePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/image_page/{page_id} [get]
// @Security ApiKeyAuth
func GetImagePage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.GetImagePage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetImagePageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		pageID, err := utils.GetURLParamInt64(r, "page_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		page, err := lpService.GetImagePage(r.Context(), &lpmodels.GetPage{
			UserID:    uID,
			PageID:    pageID,
			LessonID:  lessonID,
			ChannelID: channelID,
			PlanID:    planID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrPageNotFound):
				log.Error("image page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("image page not found"))
			default:
				log.Error("failed to get image page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("image page retrieved", slog.Int64("page_id", pageID))

		render.JSON(w, r, GetImagePageResponse{
			Response:  response.OK(),
			ImagePage: *page,
		})
	}
}

// GetVideoPage godoc
// @Summary      Get video page information
// @Description  This endpoint returns video information by ID.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Success      200 {object} pageshandler.GetVideoPageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/video_page/{page_id} [get]
// @Security ApiKeyAuth
func GetVideoPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.GetVideoPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetVideoPageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		pageID, err := utils.GetURLParamInt64(r, "page_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		page, err := lpService.GetVideoPage(r.Context(), &lpmodels.GetPage{
			UserID:    uID,
			PageID:    pageID,
			LessonID:  lessonID,
			ChannelID: channelID,
			PlanID:    planID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrPageNotFound):
				log.Error("video page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("video page not found"))
			default:
				log.Error("failed to get video page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("video page retrieved", slog.Int64("page_id", pageID))

		render.JSON(w, r, GetVideoPageResponse{
			Response:  response.OK(),
			VideoPage: *page,
		})
	}
}

// GetPDFPage godoc
// @Summary      Get pdf page information
// @Description  This endpoint returns pdf information by ID.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Success      200 {object} pageshandler.GetPDFPageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pdf_page/{page_id} [get]
// @Security ApiKeyAuth
func GetPDFPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.GetPDFPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetPDFPageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		pageID, err := utils.GetURLParamInt64(r, "page_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		page, err := lpService.GetPDFPage(r.Context(), &lpmodels.GetPage{
			UserID:    uID,
			PageID:    pageID,
			LessonID:  lessonID,
			ChannelID: channelID,
			PlanID:    planID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrPageNotFound):
				log.Error("pdf page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("pdf page not found"))
			default:
				log.Error("failed to get pdf page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("pdf page retrieved", slog.Int64("page_id", pageID))

		render.JSON(w, r, GetPDFPageResponse{
			Response: response.OK(),
			PDFPage:  *page,
		})
	}
}

// GetPages godoc
// @Summary      Get all pages from lesson relevant for user
// @Description  This endpoint returns pages information relevant for user.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param 		 limit query int false "Limit"
// @Param 		 offset query int false "Offset"
// @Success      201 {object} pageshandler.GetPagesResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pages [get]
// @Security ApiKeyAuth
func GetPages(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.GetPages"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetPagesReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		limit, err := utils.GetURLParamInt64(r, "limit")
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := utils.GetURLParamInt64(r, "offset")
		if err != nil || offset < 0 {
			offset = 0
		}

		pages, err := lpService.GetPages(r.Context(), &lpmodels.GetPages{
			UserID:    uID,
			PlanID:    planID,
			ChannelID: channelID,
			LessonID:  lessonID,
			Limit:     limit,
			Offset:    offset,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrPageNotFound):
				log.Error("pages not found", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("pages not found"))
			default:
				log.Error("failed to get pages", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("pages retrieved")

		render.JSON(w, r, GetPagesResponse{
			Response: response.OK(),
			Pages:    pages,
		})
	}
}

// UpdateImagePage godoc
// @Summary      Update image page by id
// @Description  This endpoint allows image page id and update it.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Param        pageshandler.UpdateImagePageRequest body pageshandler.UpdateImagePageRequest true "Image page updating parameters"
// @Success      200 {object} pageshandler.UpdatePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/image_page/{page_id} [patch]
// @Security ApiKeyAuth
func UpdateImagePage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.UpdateImagePage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.UpdateImagePageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		pageID, err := utils.GetURLParamInt64(r, "page_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		req, err := utils.DecodeRequestBody[UpdateImagePageRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.UpdateImagePage(r.Context(), &lpmodels.UpdateImagePage{
			UpdateBasePage: lpmodels.UpdateBasePage{
				ID:             pageID,
				ChannelID:      channelID,
				PlanID:         planID,
				LessonID:       lessonID,
				LastModifiedBy: uID,
			},
			ImageFileUrl: req.ImageFileUrl,
			ImageName:    req.ImageName,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrPageNotFound):
				log.Error("image page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("image page not found"))
			default:
				log.Error("failed to update image page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("image page updated", slog.Int64("page_id", pageID))

		render.JSON(w, r, UpdatePageResponse{
			Response:           response.OK(),
			UpdatePageResponse: *resp,
		})
	}
}

// UpdateVideoPage godoc
// @Summary      Update video page by id
// @Description  This endpoint allows video page id and update it.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Param        pageshandler.UpdateVideoPageRequest body pageshandler.UpdateVideoPageRequest true "Video page updating parameters"
// @Success      200 {object} pageshandler.UpdatePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/video_page/{page_id} [patch]
// @Security ApiKeyAuth
func UpdateVideoPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.UpdateVideoPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.UpdateVideoPageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		pageID, err := utils.GetURLParamInt64(r, "page_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		req, err := utils.DecodeRequestBody[UpdateVideoPageRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.UpdateVideoPage(r.Context(), &lpmodels.UpdateVideoPage{
			UpdateBasePage: lpmodels.UpdateBasePage{
				ID:             pageID,
				ChannelID:      channelID,
				PlanID:         planID,
				LessonID:       lessonID,
				LastModifiedBy: uID,
			},
			VideoFileUrl: req.VideoFileUrl,
			VideoName:    req.VideoName,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrPageNotFound):
				log.Error("video page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("video page not found"))
			default:
				log.Error("failed to update video page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("video page updated", slog.Int64("page_id", pageID))

		render.JSON(w, r, UpdatePageResponse{
			Response:           response.OK(),
			UpdatePageResponse: *resp,
		})
	}
}

// UpdatePDFPage godoc
// @Summary      Update pdf page by id
// @Description  This endpoint allows pdf page id and update it.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Param        pageshandler.UpdatePDFPageRequest body pageshandler.UpdatePDFPageRequest true "PDF page updating parameters"
// @Success      200 {object} pageshandler.UpdatePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pdf_page/{page_id} [patch]
// @Security ApiKeyAuth
func UpdatePDFPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.UpdatePDFPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.UpdatePDFPageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		pageID, err := utils.GetURLParamInt64(r, "page_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		req, err := utils.DecodeRequestBody[UpdatePDFPageRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.UpdatePDFPage(r.Context(), &lpmodels.UpdatePDFPage{
			UpdateBasePage: lpmodels.UpdateBasePage{
				ID:             pageID,
				ChannelID:      channelID,
				PlanID:         planID,
				LessonID:       lessonID,
				LastModifiedBy: uID,
			},
			PdfFileUrl: req.PdfFileUrl,
			PdfName:    req.PdfName,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrPageNotFound):
				log.Error("pdf page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("pdf page not found"))
			default:
				log.Error("failed to update pdf page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("pdf page updated", slog.Int64("page_id", pageID))

		render.JSON(w, r, UpdatePageResponse{
			Response:           response.OK(),
			UpdatePageResponse: *resp,
		})
	}
}

// DeletePage godoc
// @Summary      Delete page by id
// @Description  This endpoint allows page id and delete it.
// @Tags         pages
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Success      200 {object} pageshandler.DeletePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/pages/{page_id} [delete]
// @Security ApiKeyAuth
func DeletePage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.pages.DeletePage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.DeletePageReqCount.Add(r.Context(), 1)

		uID, err := utils.GetHeaderID(r, "X-User-ID")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		channelID, err := utils.GetURLParamInt64(r, "channel_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		planID, err := utils.GetURLParamInt64(r, "plan_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		lessonID, err := utils.GetURLParamInt64(r, "lesson_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}
		pageID, err := utils.GetURLParamInt64(r, "page_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		del, err := lpService.DeletePage(r.Context(), &lpmodels.DeletePage{
			UserID:    uID,
			PageID:    pageID,
			LessonID:  lessonID,
			ChannelID: channelID,
			PlanID:    planID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrPageNotFound):
				log.Error("page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("page not found"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			default:
				log.Error("failed to delete page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("page deleted", slog.Int64("page_id", pageID))

		render.JSON(w, r, DeletePageResponse{
			Response: response.OK(),
			Success:  del.Success,
		})
	}
}
