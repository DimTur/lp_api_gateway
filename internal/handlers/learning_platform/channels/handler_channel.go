package channelshandler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	resp "github.com/DimTur/lp_api_gateway/internal/lib/api/response"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidChannelID   = errors.New("invalid channel id")
	ErrChannelExitsts     = errors.New("channel already exists")
)

type CreateChannelRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=10"`
	// UserID      int64  `json:"user_id" validate:"required,numeric"`
	Public bool `json:"public" validate:"required,boolean"`
}

type CreateChannelResponce struct {
	resp.Response
	ChannelID int64 `json:"channel_id,omitempty"`
}

type GetChannelRequest struct {
	ChannelID int64 `json:"channel_id,omitempty"`
}

type GetChannelResponce struct {
	resp.Response
	Name        string `json:"name" validate:"required,name"`
	Description string `json:"description" validate:"required,description"`
	UserID      int64  `json:"user_id" validate:"required,user_id"`
	Public      bool   `json:"public" validate:"required,public"`
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

func CreateChannel(log *slog.Logger, lpService LPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.learning_platform.channels.CreateChannel"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req CreateChannelRequest

		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			log.Error("missing X-User-ID in headers")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Error("invalid X-User-ID in headers", slog.String("err", err.Error()))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			render.JSON(w, r, resp.Error("failed to decode request"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := Validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", err.Error()))
			render.JSON(w, r, resp.ValidationError(validateErr))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		respID, err := lpService.CreateChannel(
			r.Context(),
			req.Name,
			req.Description,
			userID,
			req.Public,
		)
		if errors.Is(err, ErrChannelExitsts) {
			log.Info("channel already exists", slog.String("user", req.Name))
			render.JSON(w, r, resp.Error("channel already exists"))
			w.WriteHeader(http.StatusConflict)
			return
		}
		if err != nil {
			log.Error("failed to create channel", slog.String("err", err.Error()))
			render.JSON(w, r, resp.Error("failed to create channel"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("channel created", slog.Int64("id", respID.Channel.ChannelId))

		render.JSON(w, r, CreateChannelResponce{
			Response:  resp.OK(),
			ChannelID: respID.Channel.ChannelId,
		})
		w.WriteHeader(http.StatusCreated)
	}
}
