package lessonshandler

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
	ErrLessonNotFound     = errors.New("lesson not found")
)

type LPService interface {
	CreateLesson(ctx context.Context, lesson *lpmodels.CreateLesson) (*lpmodels.CreateLessonResponse, error)
	GetLesson(ctx context.Context, lesson *lpmodels.GetLesson) (*lpmodels.GetLessonResponse, error)
	GetLessons(ctx context.Context, inputParam *lpmodels.GetLessons) ([]lpmodels.GetLessonResponse, error)
	UpdateLesson(ctx context.Context, updLesson *lpmodels.UpdateLesson) (*lpmodels.UpdateLessonResponse, error)
	DeleteLesson(ctx context.Context, delLess *lpmodels.DeleteLesson) (*lpmodels.DeleteLessonResponse, error)
}

// CreateLesson godoc
// @Summary      Create a new lesson
// @Description  This endpoint allows users to create a new lesson with the specified data.
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lessonshandler.CreateLessonRequest body lessonshandler.CreateLessonRequest true "Lesson creation parameters"
// @Success      201 {object} lessonshandler.CreateLessonResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons [post]
// @Security ApiKeyAuth
func CreateLesson(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.lessons.CreateLesson"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateChannelReqCount.Add(r.Context(), 1)

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

		req, err := utils.DecodeRequestBody[CreateLessonRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.CreateLesson(r.Context(), &lpmodels.CreateLesson{
			Name:        req.Name,
			Description: req.Description,
			CreatedBy:   uID,
			PlanID:      planID,
			ChannelID:   channelID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("invalid credentials", slog.Any("lesson", req.Name))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to create lesson", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to create lesson"))
				return
			}
		}

		log.Info("lesson created", slog.Int64("id", resp.ID))

		render.JSON(w, r, CreateLessonResponse{
			Response: response.OK(),
			LessonID: resp.ID,
		})
		w.WriteHeader(http.StatusCreated)
	}
}

// GetLesson godoc
// @Summary      Get lesson information
// @Description  This endpoint returns lesson information by ID.
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Success      200 {object} lessonshandler.GetLessonResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id} [get]
// @Security ApiKeyAuth
func GetLesson(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.lessons.GetLesson"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateChannelReqCount.Add(r.Context(), 1)

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

		lesson, err := lpService.GetLesson(r.Context(), &lpmodels.GetLesson{
			UserID:    uID,
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
				log.Error("bad request", slog.Int64("lesson_id", lessonID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrLessonNotFound):
				log.Error("lesson not found", slog.Int64("lesson_id", lessonID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("lesson not found"))
			default:
				log.Error("failed to get lesson", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("lesson retrieved", slog.Int64("lesson_id", lessonID))

		render.JSON(w, r, GetLessonResponse{
			Response: response.OK(),
			Lesson:   *lesson,
		})
	}
}

// GetLessons godoc
// @Summary      Get all lessons relevant for user
// @Description  This endpoint returns lessons information relevant for user.
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param 		 limit query int false "Limit"
// @Param 		 offset query int false "Offset"
// @Success      201 {object} lessonshandler.GetLessonsResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons [get]
// @Security ApiKeyAuth
func GetLessons(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.lessons.GetLessons"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateChannelReqCount.Add(r.Context(), 1)

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

		limit, err := utils.GetURLParamInt64(r, "limit")
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := utils.GetURLParamInt64(r, "offset")
		if err != nil || offset < 0 {
			offset = 0
		}

		lessons, err := lpService.GetLessons(r.Context(), &lpmodels.GetLessons{
			UserID:    uID,
			PlanID:    planID,
			ChannelID: channelID,
			Limit:     limit,
			Offset:    offset,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrLessonNotFound):
				log.Error("lessons not found", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("lessons not found"))
			default:
				log.Error("failed to get lessons", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("lessons retrieved")

		render.JSON(w, r, GetLessonsResponse{
			Response: response.OK(),
			Lessons:  lessons,
		})
	}
}

// UpdateLesson godoc
// @Summary      Update lesson by id
// @Description  This endpoint allows lesson id and update it.
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        lessonshandler.UpdateLessonRequest body lessonshandler.UpdateLessonRequest true "Lesson updating parameters"
// @Success      200 {object} lessonshandler.UpdateLessonResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id} [patch]
// @Security ApiKeyAuth
func UpdateLesson(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.lessons.UpdateLesson"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateChannelReqCount.Add(r.Context(), 1)

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

		req, err := utils.DecodeRequestBody[UpdateLessonRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.UpdateLesson(r.Context(), &lpmodels.UpdateLesson{
			ChannelID:      channelID,
			PlanID:         planID,
			LessonID:       lessonID,
			Name:           req.Name,
			Description:    req.Description,
			LastModifiedBy: uID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("lesson_id", lessonID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrLessonNotFound):
				log.Error("lesson not found", slog.Int64("lesson_id", lessonID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("lesson not found"))
			default:
				log.Error("failed to update lesson", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("lesson updated", slog.Int64("lesson_id", lessonID))

		render.JSON(w, r, UpdateLessonResponse{
			Response:             response.OK(),
			UpdateLessonResponse: *resp,
		})
	}
}

// DeleteLesson godoc
// @Summary      Delete lesson by id
// @Description  This endpoint allows lesson id and delete it.
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Success      200 {object} lessonshandler.DeleteLessonResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id} [delete]
// @Security ApiKeyAuth
func DeleteLesson(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.lessons.DeleteLesson"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateChannelReqCount.Add(r.Context(), 1)

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

		del, err := lpService.DeleteLesson(r.Context(), &lpmodels.DeleteLesson{
			UserID:    uID,
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
			case errors.Is(err, lpservice.ErrLessonNotFound):
				log.Error("lesson not found", slog.Int64("lesson_id", lessonID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("lesson not found"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("lesson_id", lessonID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			default:
				log.Error("failed to delete lesson", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("lesson deleted", slog.Int64("lesson_id", lessonID))

		render.JSON(w, r, DeleteLessonResponse{
			Response: response.OK(),
			Success:  del.Success,
		})
	}
}
