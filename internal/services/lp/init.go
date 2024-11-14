package lpservice

import (
	"context"
	"log/slog"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/internal/services/permissions.go"
	"github.com/go-playground/validator/v10"
)

type ChannelServiceProvider interface {
	CreateChannel(ctx context.Context, newChannel *lpmodels.CreateChannel) (*lpmodels.CreateChannelResponse, error)
	GetChannel(ctx context.Context, channel *lpmodels.GetChannel) (*lpmodels.GetChannelResponse, error)
	GetChannels(ctx context.Context, inputParam *lpmodels.GetChannelsFull) ([]lpmodels.Channel, error)
	UpdateChannel(ctx context.Context, updChannel *lpmodels.UpdateChannel) (*lpmodels.UpdateChannelResponse, error)
	DeleteChannel(ctx context.Context, delChannel *lpmodels.DelChByID) (*lpmodels.DelChByIDResp, error)
	ShareChannelToGroup(ctx context.Context, s *lpmodels.SharingChannel) (*lpmodels.SharingChannelResp, error)
	LerningGroupsShareWithChannel(ctx context.Context, channelID *lpmodels.LerningGroupsShareWithChannel) ([]string, error)
}

type LgServiceProvider interface {
	UserIsLearnerIn(ctx context.Context, user *ssomodels.UserIsLearnerIn) ([]string, error)
}

type PermissionsServiceProvider interface {
	permissions.PermissionsService
}

type LpService struct {
	Log                 *slog.Logger
	Validator           *validator.Validate
	ChannelProvider     ChannelServiceProvider
	LgServiceProvider   LgServiceProvider
	PermissionsProvider permissions.PermissionsService
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	channelProvider ChannelServiceProvider,
	lgServiceProvider LgServiceProvider,
	permissionsProvider permissions.PermissionsService,
) *LpService {
	return &LpService{
		Log:                 log,
		Validator:           validator,
		ChannelProvider:     channelProvider,
		LgServiceProvider:   lgServiceProvider,
		PermissionsProvider: permissionsProvider,
	}
}
