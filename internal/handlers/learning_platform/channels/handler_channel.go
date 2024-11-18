package channelshandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	lpservice "github.com/DimTur/lp_api_gateway/internal/services/lp"
	"github.com/DimTur/lp_api_gateway/pkg/meter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidChannelID   = errors.New("invalid channel id")
	ErrChannelExitsts     = errors.New("channel already exists")
	ErrChannelNotFound    = errors.New("channel not found")
)

type LPService interface {
	CreateChannel(ctx context.Context, newChannel *lpmodels.CreateChannel) (*lpmodels.CreateChannelResponse, error)
	GetChannel(ctx context.Context, channel *lpmodels.GetChannel) (*lpmodels.GetChannelResponse, error)
	GetChannels(ctx context.Context, inputParam *lpmodels.GetChannels) ([]lpmodels.Channel, error)
	UpdateChannel(ctx context.Context, updChannel *lpmodels.UpdateChannel) (*lpmodels.UpdateChannelResponse, error)
	DeleteChannel(ctx context.Context, delChannel *lpmodels.DelChByID) (*lpmodels.DelChByIDResp, error)
	ShareChannelToGroup(ctx context.Context, s *lpmodels.SharingChannel) (*lpmodels.SharingChannelResp, error)
}

// CreateChannel godoc
// @Summary      Create a new channel
// @Description  This endpoint allows users to create a new channel with the specified data.
// @Tags         channels
// @Accept       json
// @Produce      json
// @Param        channelshandler.CreateChannelRequest body channelshandler.CreateChannelRequest true "Channel creation parameters"
// @Success      201 {object} channelshandler.CreateChannelResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      409 {object} response.Response "Conflict"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels [post]
// @Security ApiKeyAuth
func CreateChannel(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.CreateChannel"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateChannelReqCount.Add(r.Context(), 1)

		var req CreateChannelRequest
		uID := r.Header.Get("X-User-ID")
		if uID == "" {
			log.Error("missing X-User-ID in headers")
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		resp, err := lpService.CreateChannel(r.Context(), &lpmodels.CreateChannel{
			Name:            req.Name,
			Description:     req.Description,
			CreatedBy:       uID,
			LearningGroupId: req.LearningGroupId,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("invalid credentials", slog.Any("channel", req.Name))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid credentinals"))
				return
			default:
				log.Error("failed to create channel", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("failed to create channel"))
				return
			}
		}

		log.Info("channel created", slog.Int64("id", resp.ID))

		render.JSON(w, r, CreateChannelResponse{
			Response:  response.OK(),
			ChannelID: resp.ID,
		})
		w.WriteHeader(http.StatusCreated)
	}
}

// GetChannel godoc
// @Summary      Get channel information
// @Description  This endpoint returns channel information by ID.
// @Tags         channels
// @Accept       json
// @Produce      json
// @Param        id path int true "ID of the channel"
// @Success      200 {object} channelshandler.GetChannelResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Channel not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{id} [get]
// @Security ApiKeyAuth
func GetChannel(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.GetChannel"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetChannelReqCount.Add(r.Context(), 1)

		uID := r.Header.Get("X-User-ID")
		if uID == "" {
			log.Error("missing X-User-ID in headers")
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		channelIDStr := chi.URLParam(r, "id")
		if channelIDStr == "" {
			log.Error("missing channel ID in query params")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			log.Error("invalid channel ID in query params", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		channel, err := lpService.GetChannel(r.Context(), &lpmodels.GetChannel{
			UserID:    uID,
			ChannelID: channelID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("channel_id", channelID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			case errors.Is(err, lpservice.ErrChannelNotFound):
				log.Error("channel not found", slog.Int64("channel_id", channelID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("channel not found"))
			default:
				log.Error("failed to get channel", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("channel retrieved", slog.Int64("channel_id", channelID))

		render.JSON(w, r, GetChannelResponse{
			Response: response.OK(),
			Channel:  *channel,
		})
	}
}

// GetChannels godoc
// @Summary      Get channels information
// @Description  This endpoint returns channels information relevant for user.
// @Tags         channels
// @Accept       json
// @Produce      json
// @Param 		 limit query int false "Limit"
// @Param 		 offset query int false "Offset"
// @Success      200 {object} channelshandler.GetChannelsResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Channels not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels [get]
// @Security ApiKeyAuth
func GetChannels(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.GetChannels"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetChannelReqCount.Add(r.Context(), 1)

		uID := r.Header.Get("X-User-ID")
		if uID == "" {
			log.Error("missing X-User-ID in headers")
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit <= 0 {
			limit = 10
		}

		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil || offset < 0 {
			offset = 0
		}

		channels, err := lpService.GetChannels(r.Context(), &lpmodels.GetChannels{
			UserID: uID,
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrChannelNotFound):
				log.Error("channel not found", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("channel not found"))
			default:
				log.Error("failed to get channels", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("channels retrieved")

		render.JSON(w, r, GetChannelsResponse{
			Response: response.OK(),
			Channels: channels,
		})
	}
}

// UpdateChannel godoc
// @Summary      Update channel by id
// @Description  This endpoint allows channel id and update it.
// @Tags         channels
// @Accept       json
// @Produce      json
// @Param        channelshandler.UpdateChannelRequest body channelshandler.UpdateChannelRequest true "Channels getting parameters"
// @Param        id path int true "ID of the channel"
// @Success      200 {object} channelshandler.GetChannelResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Channels not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{id} [patch]
// @Security ApiKeyAuth
func UpdateChannel(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.UpdateChannel"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetChannelReqCount.Add(r.Context(), 1)

		var req UpdateChannelRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		uID := r.Header.Get("X-User-ID")
		if uID == "" {
			log.Error("missing X-User-ID in headers")
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		channelIDStr := chi.URLParam(r, "id")
		if channelIDStr == "" {
			log.Error("missing channel ID in query params")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			log.Error("invalid channel ID in query params", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		upd, err := lpService.UpdateChannel(r.Context(), &lpmodels.UpdateChannel{
			UserID:      uID,
			ChannelID:   channelID,
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("channel_id", channelID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			default:
				log.Error("failed to update channel", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("channel updated", slog.Int64("channel_id", channelID))

		render.JSON(w, r, UpdateChannelResponse{
			Response:              response.OK(),
			UpdateChannelResponse: *upd,
		})
	}
}

// DeleteChannel godoc
// @Summary      Delete channel by id
// @Description  This endpoint allows channel id and delete it.
// @Tags         channels
// @Accept       json
// @Produce      json
// @Param        id path int true "ID of the channel"
// @Success      200 {object} channelshandler.DeleteChannelResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Channels not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{id} [delete]
// @Security ApiKeyAuth
func DeleteChannel(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.DeleteChannel"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetChannelReqCount.Add(r.Context(), 1)

		uID := r.Header.Get("X-User-ID")
		if uID == "" {
			log.Error("missing X-User-ID in headers")
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		channelIDStr := chi.URLParam(r, "id")
		if channelIDStr == "" {
			log.Error("missing channel ID in query params")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			log.Error("invalid channel ID in query params", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		del, err := lpService.DeleteChannel(r.Context(), &lpmodels.DelChByID{
			UserID:    uID,
			ChannelID: channelID,
		})
		if err != nil {
			switch {
			case errors.Is(err, lpservice.ErrPermissionDenied):
				log.Error("permissions denied", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("permissions denied"))
			case errors.Is(err, lpservice.ErrChannelNotFound):
				log.Error("channel not found", slog.Int64("channel_id", channelID))
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("channel not found"))
			case errors.Is(err, lpservice.ErrInvalidCredentials):
				log.Error("bad request", slog.Int64("channel_id", channelID))
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error("bad request"))
			default:
				log.Error("failed to delete channel", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("channel deleted", slog.Int64("channel_id", channelID))

		render.JSON(w, r, DeleteChannelResponse{
			Response: response.OK(),
			Success:  del.Success,
		})
	}
}

// ShareChannel godoc
// @Summary      Share channel by id
// @Description  This endpoint allows channel id and learning group id and share with.
// @Tags         channels
// @Accept       json
// @Produce      json
// @Param        channelshandler.ShareChannelRequest body channelshandler.ShareChannelRequest true "Channels sharing parameters"
// @Param        id path int true "ID of the channel"
// @Success      200 {object} channelshandler.ShareChannelResponse
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {object} response.Response "Unauthorized"
// @Failure      404 {object} response.Response "Channels not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /channels/{id}/share [post]
// @Security ApiKeyAuth
func ShareChannel(log *slog.Logger, val *validator.Validate, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.ShareChannel"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetChannelReqCount.Add(r.Context(), 1)

		var req ShareChannelRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		uID := r.Header.Get("X-User-ID")
		if uID == "" {
			log.Error("missing X-User-ID in headers")
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		channelIDStr := chi.URLParam(r, "id")
		if channelIDStr == "" {
			log.Error("missing channel ID in query params")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			log.Error("invalid channel ID in query params", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		s, err := lpService.ShareChannelToGroup(r.Context(), &lpmodels.SharingChannel{
			UserID:    uID,
			ChannelID: channelID,
			LGroupIDs: req.LGroupIDs,
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
				log.Error("failed to share channel", slog.String("err", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error("Internal Server Error"))
			}
		}

		log.Info("channel shared", slog.Int64("channel_id", channelID))

		render.JSON(w, r, ShareChannelResponse{
			Response: response.OK(),
			Success:  s.Success,
		})
	}
}
