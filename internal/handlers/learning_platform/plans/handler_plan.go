package planshandler

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
	ErrPlanNotFound       = errors.New("plan not found")
)

type LPService interface {
	CreatePlan(ctx context.Context, plan *lpmodels.CreatePlan) (*lpmodels.CreatePlanResponse, error)
	GetPlan(ctx context.Context, plan *lpmodels.GetPlan) (*lpmodels.GetPlanResponse, error)
	GetPlans(ctx context.Context, inputParam *lpmodels.GetPlans) ([]lpmodels.GetPlanResponse, error)
	UpdatePlan(ctx context.Context, updPlan *lpmodels.UpdatePlan) (*lpmodels.UpdatePlanResponse, error)
	DeletePlan(ctx context.Context, delPlan *lpmodels.DelPlan) (*lpmodels.DelPlanResponse, error)
	SharePlanWithUser(ctx context.Context, sharePlanWithUser *lpmodels.SharePlan) (*lpmodels.SharingPlanResp, error)
}

// CreatePlan godoc
// @Summary      Create a new plan
// @Description  This endpoint allows users to create a new plan with the specified data.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id path int true "ID of the channel"
// @Param        planshandler.CreatePlanRequest body planshandler.CreatePlanRequest true "Plan creation parameters"
// @Success      201 {object} planshandler.CreatePlanResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{id}/plans [post]
// @Security ApiKeyAuth
func CreatePlan(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.plans.CreatePlan"

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

		channelID, err := utils.GetURLParamInt64(r, "id")
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		req, err := utils.DecodeRequestBody[CreatePlanRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.CreatePlan(r.Context(), &lpmodels.CreatePlan{
			Name:            req.Name,
			Description:     req.Description,
			CreatedBy:       uID,
			ChannelID:       channelID,
			LearningGroupId: req.LearningGroupId,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("invalid credentials", slog.Any("plan", req.Name))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to create plan", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to create plan"))
				return
			}
		}

		log.Info("channel created", slog.Int64("id", resp.ID))

		render.JSON(w, r, CreatePlanResponse{
			Response: response.OK(),
			PlanID:   resp.ID,
		})
		w.WriteHeader(http.StatusCreated)
	}
}

// GetPlan godoc
// @Summary      Get plan information
// @Description  This endpoint returns plan information by ID.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Success      200 {object} planshandler.GetPlanResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Plan not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id} [get]
// @Security ApiKeyAuth
func GetPlan(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.plans.GetPlan"

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

		plan, err := lpService.GetPlan(r.Context(), &lpmodels.GetPlan{
			UserID:    uID,
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
				log.Error("bad request", slog.Int64("plan_id", planID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrPlanNotFound):
				log.Error("plan not found", slog.Int64("plan_id", planID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("plan not found"))
			default:
				log.Error("failed to get plan", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("plan retrieved", slog.Int64("plan_id", planID))

		render.JSON(w, r, GetPlanResponse{
			Response: response.OK(),
			Plan:     *plan,
		})
	}
}

// GetPlans godoc
// @Summary      Get all plans relevant for user
// @Description  This endpoint returns plans information relevant for user.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        id path int true "ID of the channel"
// @Param 		 limit query int false "Limit"
// @Param 		 offset query int false "Offset"
// @Success      201 {object} planshandler.GetPlansResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{id}/plans [get]
// @Security ApiKeyAuth
func GetPlans(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.plans.GetPlans"

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

		channelID, err := utils.GetURLParamInt64(r, "id")
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

		plans, err := lpService.GetPlans(r.Context(), &lpmodels.GetPlans{
			UserID:    uID,
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
			case errors.Is(err, lpservice.ErrPlanNotFound):
				log.Error("plan not found", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("plan not found"))
			default:
				log.Error("failed to get plan", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("plans retrieved")

		render.JSON(w, r, GetPlansResponse{
			Response: response.OK(),
			Plans:    plans,
		})
	}
}

// UpdatePlan godoc
// @Summary      Update channel by id
// @Description  This endpoint allows plan id and update it.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Param        planshandler.UpdatePlanRequest body planshandler.UpdatePlanRequest true "Plan updating parameters"
// @Success      200 {object} planshandler.UpdatePlanResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Plan not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id} [patch]
// @Security ApiKeyAuth
func UpdatePlan(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.plans.UpdatePlan"

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

		req, err := utils.DecodeRequestBody[UpdatePlanRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.UpdatePlan(r.Context(), &lpmodels.UpdatePlan{
			ChannelID:      channelID,
			PlanID:         planID,
			Name:           req.Name,
			Description:    req.Description,
			LastModifiedBy: uID,
			IsPublished:    req.IsPublished,
			Public:         req.Public,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("plan_id", planID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrPlanNotFound):
				log.Error("plan not found", slog.Int64("plan_id", planID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("plan not found"))
			default:
				log.Error("failed to update plan", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("plan updated", slog.Int64("plan_id", planID))

		render.JSON(w, r, UpdatePlanResponse{
			Response:           response.OK(),
			UpdatePlanResponse: *resp,
		})
	}
}

// DeletePlan godoc
// @Summary      Delete plan by id
// @Description  This endpoint allows plan id and delete it.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Success      200 {object} planshandler.DeletePlanResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Plan not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id} [delete]
// @Security ApiKeyAuth
func DeletePlan(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.plans.DeletePlan"

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

		del, err := lpService.DeletePlan(r.Context(), &lpmodels.DelPlan{
			UserID:    uID,
			ChannelID: channelID,
			PlanID:    planID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrPlanNotFound):
				log.Error("plan not found", slog.Int64("plan_id", planID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("plan not found"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("plan_id", planID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			default:
				log.Error("failed to delete plan", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("plan deleted", slog.Int64("plan_id", planID))

		render.JSON(w, r, DeletePlanResponse{
			Response: response.OK(),
			Success:  del.Success,
		})
	}
}

// SharePlan godoc
// @Summary      Share plan by id
// @Description  This endpoint allows plan id and user ids and share with.
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        planshandler.SharePlanRequest body planshandler.SharePlanRequest true "Plan shering parameters"
// @Param        channel_id path int true "ID of the channel"
// @Param        plan_id path int true "ID of the plan"
// @Success      200 {object} planshandler.SharePlanResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Channels not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{channel_id}/plans/{plan_id}/share [post]
// @Security ApiKeyAuth
func SharePlan(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.plans.SharePlan"

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

		req, err := utils.DecodeRequestBody[SharePlanRequest](r, log)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("bad request"))
		}

		resp, err := lpService.SharePlanWithUser(r.Context(), &lpmodels.SharePlan{
			UserID:    uID,
			ChannelID: channelID,
			PlanID:    planID,
			UsersIDs:  req.UserIDs,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			default:
				log.Error("failed to share plan", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}
		log.Info("plan shared", slog.Int64("plan_id", planID))

		render.JSON(w, r, SharePlanResponse{
			Response: response.OK(),
			Success:  resp.Success,
		})
	}
}
