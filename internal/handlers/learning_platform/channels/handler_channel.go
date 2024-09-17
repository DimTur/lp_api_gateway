package channelshandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	"github.com/DimTur/lp_api_gateway/pkg/meter"
	"github.com/DimTur/lp_api_gateway/pkg/tracer"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidChannelID   = errors.New("invalid channel id")
	ErrChannelExitsts     = errors.New("channel already exists")
	ErrChannelNotFound    = errors.New("channel not found")
)

type CreateChannelRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=10"`
	// UserID      int64  `json:"user_id" validate:"required,numeric"`
	Public bool `json:"public" validate:"required,boolean"`
}

type CreateChannelResponce struct {
	response.Response
	ChannelID int64 `json:"channel_id,omitempty"`
}

type GetChannelRequest struct {
	ChannelID int64 `json:"channel_id,omitempty"`
}

type GetChannelResponce struct {
	response.Response
	Name           string `json:"name" validate:"required,name"`
	Description    string `json:"description" validate:"required,description"`
	CreatedBy      int64  `json:"created_by" validate:"required,created_by"`
	LastModifiedBy int64  `json:"last_modified_by" validate:"required,last_modified_by"`
	Public         bool   `json:"public" validate:"required,public"`
}

type LPService interface {
	CreateChannel(
		ctx context.Context,
		name string,
		description string,
		userID int64,
		public bool,
	) (*lpv1.CreateChannelResponse, error)
	GetChannel(ctx context.Context, channelID int64) (*lpv1.GetChannelResponse, error)
}

var Validate = validator.New()

// CreateChannel godoc
// @Summary      Create a new channel
// @Description  This endpoint allows users to create a new channel with the specified data.
// @Tags         channels
// @Accept       json
// @Produce      json
// @Param        CreateChannelRequest body channelshandler.CreateChannelRequest true "Channel creation parameters"
// @Success      201 {object} channelshandler.CreateChannelResponce
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      401 {string} string "Unauthorized"
// @Failure      500 {object} response.Response "Server error"
// @Router       /create_channel [post]
// @Security ApiKeyAuth
func CreateChannel(log *slog.Logger, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.CreateChannel"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.CreateChannelReqCount.Add(r.Context(), 1)

		_, span := tracer.LPtracer.Start(r.Context(), "CreateChannel")
		defer span.End()

		var req CreateChannelRequest

		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			log.Error("missing X-User-ID in headers")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		span.AddEvent("transform_header_started")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Error("invalid X-User-ID in headers", slog.String("err", err.Error()))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		span.AddEvent("transform_completed")
		span.SetAttributes(attribute.Int64("X-User-ID", userID))

		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		span.AddEvent("validation_started")
		if err := Validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		span.AddEvent("validation_completed")
		span.SetAttributes(attribute.String("channel", req.Name))

		span.AddEvent("channel_create_started")
		respID, err := lpService.CreateChannel(
			r.Context(),
			req.Name,
			req.Description,
			userID,
			req.Public,
		)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.InvalidArgument {
					log.Error("invalid credentials", slog.Any("channel", req))
					w.WriteHeader(http.StatusBadRequest)
					render.JSON(w, r, response.Error("invalid credentials"))
					return
				}
			}

			log.Error("failed to create channel", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create channel"))
			return
		}
		span.AddEvent("channel_create_completed")
		span.SetAttributes(attribute.String("channel", req.Name))

		log.Info("channel created", slog.Int64("id", respID.Channel.ChannelId))

		render.JSON(w, r, CreateChannelResponce{
			Response:  response.OK(),
			ChannelID: respID.Channel.ChannelId,
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
// @Param        id path int true "ID канала"
// @Success      200 {object} channelshandler.GetChannelResponce
// @Failure      400 {object} response.Response "Invalid data in the request"
// @Failure      404 {object} response.Response "Channel not found"
// @Failure      500 {object} response.Response "Server error"
// @Router       /get_channel/{id} [get]
// @Security ApiKeyAuth
func GetChannel(log *slog.Logger, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.GetChannel"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		meter.AllReqCount.Add(r.Context(), 1)
		meter.GetChannelReqCount.Add(r.Context(), 1)

		_, span := tracer.LPtracer.Start(r.Context(), "GetChannel")
		defer span.End()

		channelIDStr := chi.URLParam(r, "id")
		if channelIDStr == "" {
			log.Error("missing channel ID in query params")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		span.AddEvent("transform_url_param_started")
		channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
		if err != nil {
			log.Error("invalid channel ID in query params", slog.String("err", err.Error()))
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		span.AddEvent("transform_completed")
		span.SetAttributes(attribute.Int64("channelID", channelID))

		span.AddEvent("getting_channel_started")
		channel, err := lpService.GetChannel(r.Context(), channelID)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.NotFound {
					log.Error("channel not found", slog.Int64("channel_id", channelID))
					w.WriteHeader(http.StatusNotFound)
					render.JSON(w, r, response.Error("channel does not exist"))
					return
				}
			}

			log.Error("failed to get channel", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Internal Server Error"))
			return
		}
		span.AddEvent("getting_channel_completed")
		span.SetAttributes(attribute.Int64("channelID", channelID))

		log.Info("channel retrieved", slog.Int64("channel_id", channelID))

		render.JSON(w, r, GetChannelResponce{
			Response:       response.OK(),
			Name:           channel.Channel.Name,
			Description:    channel.Channel.Description,
			CreatedBy:      channel.Channel.CreatedBy,
			LastModifiedBy: channel.Channel.LastModifiedBy,
			Public:         channel.Channel.Public,
		})
	}
}
