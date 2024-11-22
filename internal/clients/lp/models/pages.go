package lpmodels

type CreateBasePage struct {
	LessonID  int64  `json:"lesson_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
	CreatedBy string `json:"created_by" validate:"required"`
}

type CreateImagePage struct {
	CreateBasePage
	ImageFileUrl string `json:"image_file_url" validate:"required"`
	ImageName    string `json:"image_name" validate:"required"`
}

type CreateVideoPage struct {
	CreateBasePage
	VideoFileUrl string `json:"video_file_url" validate:"required"`
	VideoName    string `json:"video_name" validate:"required"`
}
type CreatePDFPage struct {
	CreateBasePage
	PdfFileUrl string `json:"pdf_file_url" validate:"required"`
	PdfName    string `json:"pdf_name" validate:"required"`
}

type CreatePageResponse struct {
	ID      int64
	Success bool
}

type GetPage struct {
	UserID    string `json:"user_id" validate:"required"`
	PageID    int64  `json:"page_id" validate:"required"`
	LessonID  int64  `json:"lesson_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
}

type BasePage struct {
	ID             int64  `json:"id"`
	LessonID       int64  `json:"lesson_id"`
	CreatedBy      string `json:"created_by"`
	LastModifiedBy string `json:"last_modified_by"`
	CreatedAt      string `json:"created_at"`
	Modified       string `json:"modified"`
	ContentType    string `json:"content_type"`
}

type ImagePage struct {
	BasePage
	ImageFileUrl string
	ImageName    string
}

type VideoPage struct {
	BasePage
	VideoFileUrl string
	VideoName    string
}

type PDFPage struct {
	BasePage
	PdfFileUrl string
	PdfName    string
}

type GetPages struct {
	UserID    string `json:"user_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
	LessonID  int64  `json:"lesson_id" validate:"required"`
	Limit     int64  `json:"limit,omitempty" validate:"min=1"`
	Offset    int64  `json:"offset,omitempty" validate:"min=0"`
}

type UpdateBasePage struct {
	ID             int64  `json:"id" validate:"required"`
	ChannelID      int64  `json:"channel_id" validate:"required"`
	PlanID         int64  `json:"plan_id" validate:"required"`
	LessonID       int64  `json:"lesson_id" validate:"required"`
	LastModifiedBy string `json:"last_modified_by" validate:"required"`
}

type UpdateImagePage struct {
	UpdateBasePage
	ImageFileUrl string `json:"image_file_url,omitempty"`
	ImageName    string `json:"image_name,omitempty"`
}

type UpdateVideoPage struct {
	UpdateBasePage
	VideoFileUrl string `json:"video_file_url,omitempty"`
	VideoName    string `json:"video_name,omitempty"`
}

type UpdatePDFPage struct {
	UpdateBasePage
	PdfFileUrl string `json:"pdf_file_url,omitempty"`
	PdfName    string `json:"pdf_name,omitempty"`
}

type UpdatePageResponse struct {
	ID      int64
	Success bool
}

type DeletePage struct {
	UserID    string `json:"user_id" validate:"required"`
	PageID    int64  `json:"page_id" validate:"required"`
	LessonID  int64  `json:"lesson_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
}

type DeletePageResponse struct {
	Success bool
}
