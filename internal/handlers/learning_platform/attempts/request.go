package attemptshandler

type UpdatePageAttemptRequest struct {
	PageID      int64  `json:"page_id" validate:"required"`
	QPAttemptID int64  `json:"question_page_attempt_id" validate:"required"`
	UserAnswer  string `json:"user_answer,omitempty"`
}
