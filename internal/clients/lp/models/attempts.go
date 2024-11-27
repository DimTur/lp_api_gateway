package lpmodels

type TryLesson struct {
	UserID    string `json:"user_id" validate:"required"`
	LessonID  int64  `json:"lesson_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
}

type QuestionPageAttempt struct {
	ID              int64  `json:"id" redis:"id"`
	PageID          int64  `json:"page_id" redis:"page_id"`
	LessonAttemptID int64  `json:"lesson_attempt_id" redis:"lesson_attempt_id"`
	IsCorrect       bool   `json:"is_correct" redis:"is_correct"`
	UserAnswer      string `json:"user_answer,omitempty"`
}

type TryLessonResp struct {
	QuestionPageAttempts []QuestionPageAttempt
}

type UpdatePageAttempt struct {
	UserID          string `json:"user_id" validate:"required"`
	LessonAttemptID int64  `json:"lesson_attempt_id" validate:"required"`
	PageID          int64  `json:"page_id" validate:"required"`
	QPAttemptID     int64  `json:"question_page_attempt_id" validate:"required"`
	UserAnswer      string `json:"user_answer,omitempty"`
}

type UpdatePageAttemptResp struct {
	Success bool
}

type CompleteLesson struct {
	UserID          string `json:"user_id" validate:"required"`
	LessonAttemptID int64  `json:"lesson_attempt_id" validate:"required"`
}

type CompleteLessonResp struct {
	ID              int64
	IsSuccessful    bool
	PercentageScore int64
}

type GetLessonAttempts struct {
	UserID   string `json:"user_id" validate:"required"`
	LessonID int64  `json:"lesson_id,omitempty"`
	Limit    int64  `json:"limit" validate:"required"`
	Offset   int64  `json:"offset" validate:"required"`
}

type LessonAttempt struct {
	ID              int64
	UserID          string
	LessonID        int64
	PlanID          int64
	ChannelID       int64
	StartTime       string
	EndTime         string
	LastModifiedBy  string
	IsComplete      bool
	IsSuccessful    bool
	PercentageScore int64
}

type GetLessonAttemptsResp struct {
	LessonAttempts []LessonAttempt
}

type LessonAttemptPermissions struct {
	UserID          string
	LessonAttemptID int64
}
