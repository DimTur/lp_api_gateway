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
	GetPlansForGroupAdmin(ctx context.Context, inputParam *lpmodels.GetPlans) ([]lpmodels.GetPlanResponse, error)
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

type PageServiceProvider interface {
	CreateImagePage(ctx context.Context, page *lpmodels.CreateImagePage) (*lpmodels.CreatePageResponse, error)
	CreateVideoPage(ctx context.Context, page *lpmodels.CreateVideoPage) (*lpmodels.CreatePageResponse, error)
	CreatePDFPage(ctx context.Context, page *lpmodels.CreatePDFPage) (*lpmodels.CreatePageResponse, error)
	GetImagePage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.ImagePage, error)
	GetVideoPage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.VideoPage, error)
	GetPDFPage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.PDFPage, error)
	GetPages(ctx context.Context, inputParams *lpmodels.GetPages) ([]lpmodels.BasePage, error)
	UpdateImagePage(ctx context.Context, updIPage *lpmodels.UpdateImagePage) (*lpmodels.UpdatePageResponse, error)
	UpdateVideoPage(ctx context.Context, updIPage *lpmodels.UpdateVideoPage) (*lpmodels.UpdatePageResponse, error)
	UpdatePDFPage(ctx context.Context, updIPage *lpmodels.UpdatePDFPage) (*lpmodels.UpdatePageResponse, error)
	DeletePage(ctx context.Context, delPage *lpmodels.DeletePage) (*lpmodels.DeletePageResponse, error)
}

type QuestionServiceProvider interface {
	CreateQuestionPage(ctx context.Context, question *lpmodels.CreateQuestionPage) (*lpmodels.CreatePageResponse, error)
	GetQuestionPage(ctx context.Context, question *lpmodels.GetPage) (*lpmodels.GetQuestionPage, error)
	UpdateQuestionPage(ctx context.Context, updQust *lpmodels.UpdateQuestionPage) (*lpmodels.UpdatePageResponse, error)
}

type AttemptServiceProvider interface {
	TryLesson(ctx context.Context, lesson *lpmodels.TryLesson) (*lpmodels.TryLessonResp, error)
	UpdatePageAttempt(ctx context.Context, attempt *lpmodels.UpdatePageAttempt) (*lpmodels.UpdatePageAttemptResp, error)
	CompleteLesson(ctx context.Context, lesson *lpmodels.CompleteLesson) (*lpmodels.CompleteLessonResp, error)
	GetLessonAttempts(ctx context.Context, inputParams *lpmodels.GetLessonAttempts) (*lpmodels.GetLessonAttemptsResp, error)
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
	PageProvider        PageServiceProvider
	QuestionProvider    QuestionServiceProvider
	AttemptProvider     AttemptServiceProvider
	LgServiceProvider   LgServiceProvider
	PermissionsProvider permissions.PermissionsService
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	channelProvider ChannelServiceProvider,
	planProvider PlanServiceProvider,
	lessonProvider LessonServiceProvider,
	pageProvider PageServiceProvider,
	questionProvider QuestionServiceProvider,
	attemptProvider AttemptServiceProvider,
	lgServiceProvider LgServiceProvider,
	permissionsProvider permissions.PermissionsService,
) *LpService {
	return &LpService{
		Log:                 log,
		Validator:           validator,
		ChannelProvider:     channelProvider,
		PlanProvider:        planProvider,
		LessonProvider:      lessonProvider,
		PageProvider:        pageProvider,
		QuestionProvider:    questionProvider,
		AttemptProvider:     attemptProvider,
		LgServiceProvider:   lgServiceProvider,
		PermissionsProvider: permissionsProvider,
	}
}
