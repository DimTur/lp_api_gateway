package attemptshandler

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
	ErrLessonAttemtNotFound       = errors.New("lesson attempt not found")
	ErrQuestionPageAttemtNotFound = errors.New("question page attempt not found")
	ErrAnswerNotFound             = errors.New("page answer not found")
	ErrPermissionsDenied          = errors.New("permissions denied")
)

type LPService interface {
	TryLesson(ctx context.Context, lesson *lpmodels.TryLesson) (*lpmodels.TryLessonResp, error)
	UpdatePageAttempt(ctx context.Context, attempt *lpmodels.UpdatePageAttempt) (*lpmodels.UpdatePageAttemptResp, error)
	CompleteLesson(ctx context.Context, lesson *lpmodels.CompleteLesson) (*lpmodels.CompleteLessonResp, error)
	GetLessonAttempts(ctx context.Context, inputParams *lpmodels.GetLessonAttempts) (*lpmodels.GetLessonAttemptsResp, error)
}

// TryLesson godoc
// @Summary      Create a new lesson attempt or get it if exist not completed
// @Description  This endpoint allows user, channel, plan, lesson id and create a new lesson attempt with the specified data.
// @Tags         attempts
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Success      201 {object} attemptshandler.TryLessonResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/attempts [post]
// @Security ApiKeyAuth
func TryLesson(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.attempts.TryLesson"

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

		resp, err := lpService.TryLesson(r.Context(), &lpmodels.TryLesson{
			UserID:    uID,
			LessonID:  lessonID,
			PlanID:    planID,
			ChannelID: channelID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("invalid credentials", slog.Any("lesson_id", lessonID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			case errors.Is(err, lpservice.ErrQuestionPageAttemtNotFound):
				log.Error("question page attempt not found", slog.Any("lesson_id", lessonID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("question page attempt not found"))
				return
			default:
				log.Error("failed to get or create lesson attempt", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to get or create lesson attempt"))
				return
			}
		}

		log.Info("lesson attemp got")

		render.JSON(w, r, TryLessonResponse{
			Response:             response.OK(),
			QuestionPageAttempts: resp.QuestionPageAttempts,
		})
		w.WriteHeader(http.StatusCreated)
	}
}

// UpdatePageAttempt godoc
// @Summary      Update question page attempt by id
// @Description  This endpoint allows question page attempt id and update it.
// @Tags         attempts
// @Accept       json
// @Produce      json
// @Param        lesson_attempt_id path int true "ID of the lesson attempt"
// @Param        attemptshandler.UpdatePageAttemptRequest body attemptshandler.UpdatePageAttemptRequest true "Question page attempt updating parameters"
// @Success      200 {object} attemptshandler.UpdatePageAttemptResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Question page attempt not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /lessons/attempts/{lesson_attempt_id} [patch]
// @Security ApiKeyAuth
func UpdatePageAttempt(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.attempts.UpdatePageAttempt"

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

		lessonAttemptID, err := utils.GetURLParamInt64(r, "lesson_attempt_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		req, err := utils.DecodeRequestBody[UpdatePageAttemptRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.UpdatePageAttempt(r.Context(), &lpmodels.UpdatePageAttempt{
			UserID:          uID,
			LessonAttemptID: lessonAttemptID,
			PageID:          req.PageID,
			QPAttemptID:     req.QPAttemptID,
			UserAnswer:      req.UserAnswer,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("question_page_attempt_id", req.QPAttemptID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrAnswerNotFound):
				log.Error("question page attempt not found", slog.Int64("question_page_attempt_id", req.QPAttemptID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("question page attempt not found"))
			default:
				log.Error("failed to update question page attempt", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("question page attempt updated", slog.Int64("question_page_attempt_id", req.QPAttemptID))

		render.JSON(w, r, UpdatePageAttemptResponse{
			Response: response.OK(),
			Success:  resp.Success,
		})
	}
}

// CompleteLesson godoc
// @Summary      Complete lesson attempt by id
// @Description  This endpoint allows lesson attempt id and update it.
// @Tags         attempts
// @Accept       json
// @Produce      json
// @Param        lesson_attempt_id path int true "ID of the lesson attempt"
// @Success      200 {object} attemptshandler.UpdatePageAttemptResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Question page attempt not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /lessons/attempts/{lesson_attempt_id}/complete [patch]
// @Security ApiKeyAuth
func CompleteLesson(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.attempts.CompleteLesson"

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

		lessonAttemptID, err := utils.GetURLParamInt64(r, "lesson_attempt_id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.CompleteLesson(r.Context(), &lpmodels.CompleteLesson{
			UserID:          uID,
			LessonAttemptID: lessonAttemptID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("lesson_attempt_id", lessonAttemptID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrQuestionPageAttemtNotFound):
				log.Error("lesson attempt not found", slog.Int64("lesson_attempt_id", lessonAttemptID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("lesson attempt not found"))
			default:
				log.Error("failed to complete lesson attempt", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("lesson attempt completed", slog.Int64("lesson_attempt_id", lessonAttemptID))

		render.JSON(w, r, CompleteLessonResponse{
			Response:        response.OK(),
			ID:              resp.ID,
			IsSuccessful:    resp.IsSuccessful,
			PercentageScore: resp.PercentageScore,
		})
	}
}

// GetLessonAttempts godoc
// @Summary      Get all lesson attempts relevant for user
// @Description  This endpoint returns lesson attempts information relevant for user.
// @Tags         attempts
// @Accept       json
// @Produce      json
// @Param        lesson_id path int true "ID of the lesson"
// @Param 		 limit query int false "Limit"
// @Param 		 offset query int false "Offset"
// @Success      201 {object} attemptshandler.LessonAttemptsResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson attempts not found"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /lessons/{lesson_id}/attempts [get]
// @Security ApiKeyAuth
func GetLessonAttempts(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.attempts.GetLessonAttempts"

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

		attempts, err := lpService.GetLessonAttempts(r.Context(), &lpmodels.GetLessonAttempts{
			UserID:   uID,
			LessonID: lessonID,
			Limit:    limit,
			Offset:   offset,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrLessonAttemtNotFound):
				log.Error("lesson attempt not found", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("lesson attempt not found"))
			default:
				log.Error("failed to get lesson attempt", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("lesson attempts retrieved")

		render.JSON(w, r, LessonAttemptsResponse{
			Response:       response.OK(),
			LessonAttempts: attempts.LessonAttempts,
		})
	}
}
