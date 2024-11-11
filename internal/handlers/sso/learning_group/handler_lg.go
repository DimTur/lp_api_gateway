package learninggrouphandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	ssoservice "github.com/DimTur/lp_api_gateway/internal/services/sso"
	"github.com/DimTur/lp_api_gateway/pkg/meter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidGroupID     = errors.New("invalid group id")
	ErrGroupExists        = errors.New("group already exists")
	ErrGroupNotFound      = errors.New("group not found")
	ErrPermissionDenied   = errors.New("you don't have permissions")

	ErrInternal = errors.New("internal error")
)

type LgService interface {
	CreateLearningGroup(ctx context.Context, newLg *ssomodels.CreateLearningGroup) (*ssomodels.CreateLearningGroupResp, error)
	GetLearningGroupByID(ctx context.Context, lgID *ssomodels.GetLgByID) (*ssomodels.GetLgByIDResp, error)
	UpdateLearningGroup(ctx context.Context, updFields *ssomodels.UpdateLearningGroup) (*ssomodels.UpdateLearningGroupResp, error)
	DeleteLearningGroup(ctx context.Context, lgID *ssomodels.DelLgByID) (*ssomodels.DelLgByIDResp, error)
	GetLearningGroups(ctx context.Context, uID *ssomodels.GetLGroups) (*ssomodels.GetLGroupsResp, error)
}

// CreateLearningGroup godoc
// @Summary      Create a new learning group
// @Description  This endpoint allows learning group info and create it.
// @Tags         learning groups
// @Accept       json
// @Produce      json
// @Param        learninggrouphandler.CreateLearningGroupRequest body learninggrouphandler.CreateLearningGroupRequest true "Creating parameters"
// @Success      201 {object} learninggrouphandler.CreateLGroupResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /learning_groups [post]
// @Security ApiKeyAuth
func CreateLearningGroup(log *slog.Logger, val *validator.Validate, lgService LgService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sso.learning_group.CreateLearningGroup"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignUpReqCount.Add(r.Context(), 1)

		uID := r.Header.Get("X-User-ID")
		var req CreateLearningGroupRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request from", r.Header.Get("X-User-ID")))

		resp, err := lgService.CreateLearningGroup(r.Context(), &ssomodels.CreateLearningGroup{
			Name:        req.Name,
			CreatedBy:   uID,
			ModifiedBy:  uID,
			GroupAdmins: []string{uID},
			Learners:    []string{uID},
		})
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
				log.Error("invalid credentinals", slog.Any("name", req.Name))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			case errors.Is(err, ssoservice.ErrGroupExists):
				log.Error("learning group already exists", slog.Any("name", req.Name))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("learning group already exists"))
				return
			default:
				log.Error("creating learning group failed", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("creating learning group failed"))
				return
			}
		}

		log.Info("learning group created successfully")

		render.JSON(w, r, CreateLGroupResponse{
			Response: response.OK(),
			Success:  resp.Success,
		})
	}
}

// GetLearningGroupByID godoc
// @Summary      Get learning group by id
// @Description  This endpoint allows learning group id and returns lg info.
// @Tags         learning groups
// @Accept       json
// @Produce      json
// @Param        id path string true "ID of the learning group"
// @Success      200 {object} learninggrouphandler.GetLgByIDResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      404 {object} response.Response "Not Found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /learning_group/{id} [get]
// @Security ApiKeyAuth
func GetLearningGroupByID(log *slog.Logger, val *validator.Validate, lgService LgService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sso.learning_group.CreateLearningGroup"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignUpReqCount.Add(r.Context(), 1)

		uID := r.Header.Get("X-User-ID")
		lgID := chi.URLParam(r, "id")
		if lgID == "" {
			log.Error("missing learning group ID in query params")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		log.Info("request received to get learning group",
			slog.Any("request from", uID),
			slog.String("learning group id", lgID),
		)

		resp, err := lgService.GetLearningGroupByID(r.Context(), &ssomodels.GetLgByID{
			UserID: uID,
			LgId:   lgID,
		})
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
				log.Error("invalid credentinals", slog.Any("learning_group_id", lgID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			case errors.Is(err, ssoservice.ErrGroupNotFound):
				log.Error("learning group not found", slog.Any("learning_group_id", lgID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("learning group not found"))
				return
			default:
				log.Error("failed to get learning group", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to get learning group"))
				return
			}
		}

		log.Info("learning group got successfully")

		render.JSON(w, r, GetLgByIDResponse{
			Response:      response.OK(),
			LearningGroup: resp,
		})
	}
}

// UpdateLearningGroup godoc
// @Summary      Update learning group by id
// @Description  This endpoint allows learning group id and update it.
// @Tags         learning groups
// @Accept       json
// @Produce      json
// @Param        id path string true "ID of the learning group"
// @Param        learninggrouphandler.UpdateLearningGroupRequest body learninggrouphandler.UpdateLearningGroupRequest true "Getting parameters"
// @Success      200 {object} learninggrouphandler.UpdateLGroupResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      404 {object} response.Response "Not Found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /learning_group/{id} [patch]
// @Security ApiKeyAuth
func UpdateLearningGroup(log *slog.Logger, val *validator.Validate, lgService LgService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sso.learning_group.CreateLearningGroup"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignUpReqCount.Add(r.Context(), 1)

		var req UpdateLearningGroupRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		uID := r.Header.Get("X-User-ID")
		lgID := chi.URLParam(r, "id")
		if lgID == "" {
			log.Error("missing learning group ID in query params")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		log.Info("request received to update learning group",
			slog.Any("request from", uID),
			slog.String("learning group id", lgID),
		)

		resp, err := lgService.UpdateLearningGroup(r.Context(), &ssomodels.UpdateLearningGroup{
			UserID:      uID,
			LgId:        lgID,
			Name:        req.Name,
			ModifiedBy:  uID,
			GroupAdmins: req.GroupAdmins,
			Learners:    req.Learners,
		})
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrGroupNotFound):
				log.Error("learning group not found", slog.Any("learning_group_id", lgID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("learning group not found"))
				return
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
				log.Error("invalid credentinals", slog.Any("learning_group_id", lgID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to update learning group", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to update learning group"))
				return
			}
		}

		log.Info("learning group got successfully")

		render.JSON(w, r, UpdateLGroupResponse{
			Response: response.OK(),
			Success:  resp.Success,
		})
	}
}

// DeleteLearningGroup godoc
// @Summary      Delete learning group by id
// @Description  This endpoint allows learning group id and delete it.
// @Tags         learning groups
// @Accept       json
// @Produce      json
// @Param        id path string true "ID of the learning group"
// @Success      200 {object} learninggrouphandler.DeleteLGroupResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      500 {object} response.Response "Server error"
// @Router       /learning_group/{id} [delete]
// @Security ApiKeyAuth
func DeleteLearningGroup(log *slog.Logger, val *validator.Validate, lgService LgService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sso.learning_group.DeleteLearningGroup"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignUpReqCount.Add(r.Context(), 1)

		uID := r.Header.Get("X-User-ID")
		lgID := chi.URLParam(r, "id")
		if lgID == "" {
			log.Error("missing learning group ID in query params")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		log.Info("request received to delete learning group",
			slog.Any("request from", uID),
			slog.String("learning group id", lgID),
		)

		resp, err := lgService.DeleteLearningGroup(r.Context(), &ssomodels.DelLgByID{
			UserID: uID,
			LgID:   lgID,
		})
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.Any("learning_group_id", lgID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
				return
			case errors.Is(err, ssoservice.ErrInvalidCredentials):
				log.Error("invalid credentinals", slog.Any("learning_group_id", lgID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to delete learning group", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to delete learning group"))
				return
			}
		}

		log.Info("learning group deleted successfully")

		render.JSON(w, r, DeleteLGroupResponse{
			Response: response.OK(),
			Success:  resp.Success,
		})
	}
}

// GetLearningGroups godoc
// @Summary      Get learning groups relevant for user
// @Description  This endpoint allows user id and returns all relevant learning groups.
// @Tags         learning groups
// @Accept       json
// @Produce      json
// @Success      200 {object} learninggrouphandler.GetLearningGroupsResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      404 {object} response.Response "Not Found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /learning_groups [get]
// @Security ApiKeyAuth
func GetLearningGroups(log *slog.Logger, val *validator.Validate, lgService LgService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.sso.learning_group.GetLearningGroups"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.SignUpReqCount.Add(r.Context(), 1)

		uID := r.Header.Get("X-User-ID")
		log.Info("request received to get learning group", slog.Any("request from", uID))

		resp, err := lgService.GetLearningGroups(r.Context(), &ssomodels.GetLGroups{
			UserID: uID,
		})
		if err != nil {
			switch {
			case errors.Is(err, ssoservice.ErrGroupNotFound):
				log.Error("learning group not found", slog.Any("user_id", uID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("learning groups not found"))
				return
			default:
				log.Error("failed to get learning group", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to get learning group"))
				return
			}
		}

		log.Info("learning group got successfully")

		render.JSON(w, r, GetLearningGroupsResponse{
			Response:       response.OK(),
			LearningGroups: resp,
		})
	}
}
