package permissions

import (
	"context"
	"log/slog"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/go-playground/validator/v10"
)

type ChannelPermissionsProvider interface {
	IsChannelCreator(ctx context.Context, isCC *lpmodels.IsChannelCreator) (*lpmodels.IsChannelCreatorResp, error)
	LerningGroupsShareWithChannel(ctx context.Context, channelID *lpmodels.LerningGroupsShareWithChannel) ([]string, error)
}

type LgPermissionsProvider interface {
	UserIsGroupAdminIn(ctx context.Context, user *ssomodels.UserIsGroupAdminIn) ([]string, error)
	UserIsLearnerIn(ctx context.Context, user *ssomodels.UserIsLearnerIn) ([]string, error)
	IsGroupAdmin(ctx context.Context, uIsGroupAdmin *ssomodels.IsGroupAdmin) (*ssomodels.IsGroupAdminResp, error)
}

type RedisPermissionsProvider interface {
	SaveLgUser(ctx context.Context, userID string, groupIDs []string) error
	SaveLgShareWithChannel(ctx context.Context, channelID int64, groupIDs []string) error
	CheckGroupsIntersection(ctx context.Context, userID string, channelID int64) (bool, error)
}

type PermissionsService struct {
	log                        *slog.Logger
	validator                  *validator.Validate
	channelPermissionsProvider ChannelPermissionsProvider
	lgPermissionsProvider      LgPermissionsProvider
	redisPermissionsProvider   RedisPermissionsProvider
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	channelPermissionsProvider ChannelPermissionsProvider,
	lgPermissionsProvider LgPermissionsProvider,
	redisPermissionsProvider RedisPermissionsProvider,
) *PermissionsService {
	return &PermissionsService{
		log:                        log,
		validator:                  validator,
		channelPermissionsProvider: channelPermissionsProvider,
		lgPermissionsProvider:      lgPermissionsProvider,
		redisPermissionsProvider:   redisPermissionsProvider,
	}
}
