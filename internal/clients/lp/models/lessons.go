package lpmodels

type CreateLesson struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	CreatedBy   string `json:"created_by" validate:"required"`
	PlanID      int64  `json:"plan_id" validate:"required"`
	ChannelID   int64  `json:"channel_id" validate:"required"`
}

type CreateLessonResponse struct {
	ID      int64
	Success bool
}

type GetLesson struct {
	UserID    string `json:"user_id" validate:"required"`
	LessonID  int64  `json:"lesson_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
}

type GetLessonResponse struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      string
	LastModifiedBy string
	CreatedAt      string
	Modified       string
}

type GetLessons struct {
	UserID    string `json:"user_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
	Limit     int64  `json:"limit,omitempty" validate:"min=1"`
	Offset    int64  `json:"offset,omitempty" validate:"min=0"`
}

type UpdateLesson struct {
	ChannelID      int64  `json:"channel_id" validate:"required"`
	PlanID         int64  `json:"plan_id" validate:"required"`
	LessonID       int64  `json:"lesson_id" validate:"required"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	LastModifiedBy string `json:"last_modified_by" validate:"required"`
}

type UpdateLessonResponse struct {
	ID      int64
	Success bool
}

type DeleteLesson struct {
	UserID    string `json:"user_id" validate:"required"`
	LessonID  int64  `json:"lesson_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
}

type DeleteLessonResponse struct {
	Success bool
}
