package meter

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	metr "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	// Meter for request monitoring
	ReqMeter = otel.Meter("requests-meter")

	// Global counters
	AllReqCount, _ = ReqMeter.Int64Counter("requests_total", metr.WithDescription("Total number of requests"))

	// Auth
	SignUpReqCount, _     = ReqMeter.Int64Counter("requests_sign_up", metr.WithDescription("Sign up number of requests"))
	SignInReqCount, _     = ReqMeter.Int64Counter("requests_sign_in", metr.WithDescription("Sign in number of requests"))
	SignInByTgReqCount, _ = ReqMeter.Int64Counter("requests_sign_in_by_tg", metr.WithDescription("Sign in by Telegram number of requests"))
	CheckOtpReqCount, _   = ReqMeter.Int64Counter("requests_check_otp", metr.WithDescription("Check OTP number of requests"))
	UpdateInfoReqCount, _ = ReqMeter.Int64Counter("requests_update_info", metr.WithDescription("Update user info number of requests"))

	// Learning Groups
	CreateLearningGroupReqCount, _ = ReqMeter.Int64Counter("requests_create_learning_group", metr.WithDescription("Create Learning Group number of requests"))
	GetLearningGroupReqCount, _    = ReqMeter.Int64Counter("requests_get_learning_group", metr.WithDescription("Get Learning Group by ID number of requests"))
	UpdateLearningGroupReqCount, _ = ReqMeter.Int64Counter("requests_update_learning_group", metr.WithDescription("Update Learning Group number of requests"))
	DeleteLearningGroupReqCount, _ = ReqMeter.Int64Counter("requests_delete_learning_group", metr.WithDescription("Delete Learning Group number of requests"))
	GetLearningGroupsReqCount, _   = ReqMeter.Int64Counter("requests_get_learning_groups", metr.WithDescription("Get all Learning Groups number of requests"))

	// Channels
	CreateChannelReqCount, _ = ReqMeter.Int64Counter("requests_create_channel", metr.WithDescription("Create Channel number of requests"))
	GetChannelReqCount, _    = ReqMeter.Int64Counter("requests_get_channel", metr.WithDescription("Get Channel by ID number of requests"))
	UpdateChannelReqCount, _ = ReqMeter.Int64Counter("requests_update_channel", metr.WithDescription("Update Channel number of requests"))
	DeleteChannelReqCount, _ = ReqMeter.Int64Counter("requests_delete_channel", metr.WithDescription("Delete Channel number of requests"))
	ShareChannelReqCount, _  = ReqMeter.Int64Counter("requests_share_channel", metr.WithDescription("Share Channel number of requests"))
	GetChannelsReqCount, _   = ReqMeter.Int64Counter("requests_get_channels", metr.WithDescription("Get all Channels number of requests"))

	// Plans
	CreatePlanReqCount, _ = ReqMeter.Int64Counter("requests_create_plan", metr.WithDescription("Create Plan number of requests"))
	GetPlanReqCount, _    = ReqMeter.Int64Counter("requests_get_plan", metr.WithDescription("Get Plan by ID number of requests"))
	UpdatePlanReqCount, _ = ReqMeter.Int64Counter("requests_update_plan", metr.WithDescription("Update Plan number of requests"))
	DeletePlanReqCount, _ = ReqMeter.Int64Counter("requests_delete_plan", metr.WithDescription("Delete Plan number of requests"))
	SharePlanReqCount, _  = ReqMeter.Int64Counter("requests_share_plan", metr.WithDescription("Share Plan number of requests"))
	GetPlansReqCount, _   = ReqMeter.Int64Counter("requests_get_plans", metr.WithDescription("Get all Plans number of requests"))

	// Lessons
	CreateLessonReqCount, _ = ReqMeter.Int64Counter("requests_create_lesson", metr.WithDescription("Create Lesson number of requests"))
	GetLessonReqCount, _    = ReqMeter.Int64Counter("requests_get_lesson", metr.WithDescription("Get Lesson by ID number of requests"))
	UpdateLessonReqCount, _ = ReqMeter.Int64Counter("requests_update_lesson", metr.WithDescription("Update Lesson number of requests"))
	DeleteLessonReqCount, _ = ReqMeter.Int64Counter("requests_delete_lesson", metr.WithDescription("Delete Lesson number of requests"))
	GetLessonsReqCount, _   = ReqMeter.Int64Counter("requests_get_lessons", metr.WithDescription("Get all Lessons number of requests"))

	// Pages
	CreateImagePageReqCount, _ = ReqMeter.Int64Counter("requests_create_image_page", metr.WithDescription("Create Image Page number of requests"))
	CreateVideoPageReqCount, _ = ReqMeter.Int64Counter("requests_create_video_page", metr.WithDescription("Create Video Page number of requests"))
	CreatePDFPageReqCount, _   = ReqMeter.Int64Counter("requests_create_pdf_page", metr.WithDescription("Create PDF Page number of requests"))
	GetImagePageReqCount, _    = ReqMeter.Int64Counter("requests_get_image_page", metr.WithDescription("Get Image Page by ID number of requests"))
	GetVideoPageReqCount, _    = ReqMeter.Int64Counter("requests_get_video_page", metr.WithDescription("Get Video Page by ID number of requests"))
	GetPDFPageReqCount, _      = ReqMeter.Int64Counter("requests_get_pdf_page", metr.WithDescription("Get PDF Page by ID number of requests"))
	UpdateImagePageReqCount, _ = ReqMeter.Int64Counter("requests_update_image_page", metr.WithDescription("Update Image Page number of requests"))
	UpdateVideoPageReqCount, _ = ReqMeter.Int64Counter("requests_update_video_page", metr.WithDescription("Update Video Page number of requests"))
	UpdatePDFPageReqCount, _   = ReqMeter.Int64Counter("requests_update_pdf_page", metr.WithDescription("Update PDF Page number of requests"))
	DeletePageReqCount, _      = ReqMeter.Int64Counter("requests_delete_page", metr.WithDescription("Delete Page number of requests"))
	GetPagesReqCount, _        = ReqMeter.Int64Counter("requests_get_pages", metr.WithDescription("Get all Pages number of requests"))

	// Questions
	CreateQuestionPageReqCount, _ = ReqMeter.Int64Counter("requests_create_question_page", metr.WithDescription("Create Question Page number of requests"))
	GetQuestionPageReqCount, _    = ReqMeter.Int64Counter("requests_get_question_page", metr.WithDescription("Get Question Page by ID number of requests"))
	UpdateQuestionPageReqCount, _ = ReqMeter.Int64Counter("requests_update_question_page", metr.WithDescription("Update Question Page number of requests"))

	// Attempts
	TryLessonReqCount, _         = ReqMeter.Int64Counter("requests_try_lesson", metr.WithDescription("Try Lesson number of requests"))
	UpdatePageAttemptReqCount, _ = ReqMeter.Int64Counter("requests_update_page_attempt", metr.WithDescription("Update Page Attempt number of requests"))
	CompleteLessonReqCount, _    = ReqMeter.Int64Counter("requests_complete_lesson", metr.WithDescription("Complete Lesson number of requests"))
	GetLessonAttemptsReqCount, _ = ReqMeter.Int64Counter("requests_get_lesson_attempts", metr.WithDescription("Get Lesson Attempts number of requests"))
)

func InitMeter(ctx context.Context, serviceName string) (*metric.MeterProvider, error) {
	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String("service.name", serviceName),
	))
	if err != nil {
		return nil, err
	}

	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}
