package questionshandler

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

type LPService interface {
	CreateQuestionPage(ctx context.Context, question *lpmodels.CreateQuestionPage) (*lpmodels.CreatePageResponse, error)
	GetQuestionPage(ctx context.Context, question *lpmodels.GetPage) (*lpmodels.GetQuestionPage, error)
	UpdateQuestionPage(ctx context.Context, updQust *lpmodels.UpdateQuestionPage) (*lpmodels.UpdatePageResponse, error)
}

// CreateQuestionPage godoc
// @Summary      Create a new question page
// @Description  This endpoint allows users to create a new questin page with the specified data.
// @Tags         questions
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        questionshandler.CreateQuestionPageRequest body questionshandler.CreateQuestionPageRequest true "Question page creation parameters"
// @Success      201 {object} questionshandler.CreatePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/question_page [post]
// @Security ApiKeyAuth
func CreateQuestionPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.questions.CreateQuestionPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateQuestionPageReqCount.Add(r.Context(), 1)

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

		req, err := utils.DecodeRequestBody[CreateQuestionPageRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.CreateQuestionPage(r.Context(), &lpmodels.CreateQuestionPage{
			LessonID:  lessonID,
			PlanID:    planID,
			ChannelID: channelID,
			CreatedBy: uID,
			Question:  req.Question,
			OptionA:   req.OptionA,
			OptionB:   req.OptionB,
			OptionC:   req.OptionC,
			OptionD:   req.OptionD,
			OptionE:   req.OptionE,
			Answer:    req.Answer,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("invalid credentials", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to create question page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to create question page"))
				return
			}
		}

		log.Info("question page created", slog.Int64("id", resp.ID))

		render.JSON(w, r, CreatePageResponse{
			Response: response.OK(),
			PageID:   resp.ID,
		})
		w.WriteHeader(http.StatusCreated)
	}
}

// GetQuestionPage godoc
// @Summary      Get question page information
// @Description  This endpoint returns question information by ID.
// @Tags         questions
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Success      200 {object} questionshandler.GetQuestionPageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Lesson not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/question_page/{page_id} [get]
// @Security ApiKeyAuth
func GetQuestionPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.questions.GetQuestionPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetQuestionPageReqCount.Add(r.Context(), 1)

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

		page, err := lpService.GetQuestionPage(r.Context(), &lpmodels.GetPage{
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
			case errors.Is(err, lpservice.ErrQuestionNotFound):
				log.Error("question page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("question page not found"))
			default:
				log.Error("failed to get question page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("question page retrieved", slog.Int64("page_id", pageID))

		render.JSON(w, r, GetQuestionPageResponse{
			Response:     response.OK(),
			QuestionPage: *page,
		})
	}
}

// UpdateQuestionPage godoc
// @Summary      Update question page by id
// @Description  This endpoint allows question page id and update it.
// @Tags         questions
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        lesson_id path int true "ID of the lesson"
// @Param        page_id path int true "ID of the page"
// @Param        questionshandler.UpdateQuestinPageRequest body questionshandler.UpdateQuestinPageRequest true "Question page updating parameters"
// @Success      200 {object} questionshandler.UpdatePageResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Question not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/lessons/{lesson_id}/question_page/{page_id} [patch]
// @Security ApiKeyAuth
func UpdateQuestionPage(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.questions.UpdateQuestionPage"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.UpdateQuestionPageReqCount.Add(r.Context(), 1)

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

		req, err := utils.DecodeRequestBody[UpdateQuestinPageRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.UpdateQuestionPage(r.Context(), &lpmodels.UpdateQuestionPage{
			ID:             pageID,
			ChannelID:      channelID,
			PlanID:         planID,
			LessonID:       lessonID,
			LastModifiedBy: uID,
			Question:       req.Question,
			OptionA:        req.OptionA,
			OptionB:        req.OptionB,
			OptionC:        req.OptionC,
			OptionD:        req.OptionD,
			OptionE:        req.OptionE,
			Answer:         req.Answer,
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
			case errors.Is(err, lpservice.ErrQuestionNotFound):
				log.Error("question page not found", slog.Int64("page_id", pageID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("question page not found"))
			default:
				log.Error("failed to update question page", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("question page updated", slog.Int64("page_id", pageID))

		render.JSON(w, r, UpdatePageResponse{
			Response:           response.OK(),
			UpdatePageResponse: *resp,
		})
	}
}
