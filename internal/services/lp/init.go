package lpservice

import (
	"context"
	"log/slog"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/internal/services/permissions"
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

type PlanServiceProvider interface {
	CreatePlan(ctx context.Context, plan *lpmodels.CreatePlan) (*lpmodels.CreatePlanResponse, error)
	GetPlan(ctx context.Context, plan *lpmodels.GetPlan) (*lpmodels.GetPlanResponse, error)
	GetPlans(ctx context.Context, inputParam *lpmodels.GetPlans) ([]lpmodels.GetPlanResponse, error)
	UpdatePlan(ctx context.Context, updPlan *lpmodels.UpdatePlan) (*lpmodels.UpdatePlanResponse, error)
	DeletePlan(ctx context.Context, delPlan *lpmodels.DelPlan) (*lpmodels.DelPlanResponse, error)
	SharePlanWithUser(ctx context.Context, sharePlanWithUser *lpmodels.SharePlan) (*lpmodels.SharingPlanResp, error)
}

type LessonServiceProvider interface {
	CreateLesson(ctx context.Context, lesson *lpmodels.CreateLesson) (*lpmodels.CreateLessonResponse, error)
	GetLesson(ctx context.Context, lesson *lpmodels.GetLesson) (*lpmodels.GetLessonResponse, error)
	GetLessons(ctx context.Context, inputParam *lpmodels.GetLessons) ([]lpmodels.GetLessonResponse, error)
	UpdateLesson(ctx context.Context, updLesson *lpmodels.UpdateLesson) (*lpmodels.UpdateLessonResponse, error)
	DeleteLesson(ctx context.Context, delLess *lpmodels.DeleteLesson) (*lpmodels.DeleteLessonResponse, error)
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
	PlanProvider        PlanServiceProvider
	LessonProvider      LessonServiceProvider
	LgServiceProvider   LgServiceProvider
	PermissionsProvider permissions.PermissionsService
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	channelProvider ChannelServiceProvider,
	planProvider PlanServiceProvider,
	lessonProvider LessonServiceProvider,
	lgServiceProvider LgServiceProvider,
	permissionsProvider permissions.PermissionsService,
) *LpService {
	return &LpService{
		Log:                 log,
		Validator:           validator,
		ChannelProvider:     channelProvider,
		PlanProvider:        planProvider,
		LessonProvider:      lessonProvider,
		LgServiceProvider:   lgServiceProvider,
		PermissionsProvider: permissionsProvider,
	}
}