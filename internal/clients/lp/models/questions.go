package lpmodels

type CreateQuestionPage struct {
	LessonID  int64  `json:"lesson_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
	CreatedBy string `json:"created_by" validate:"required"`

	Question string `json:"question" validate:"required"`
	OptionA  string `json:"option_a" validate:"required"`
	OptionB  string `json:"option_b" validate:"required"`
	OptionC  string `json:"option_c,omitempty"`
	OptionD  string `json:"option_d,omitempty"`
	OptionE  string `json:"option_e,omitempty"`
	Answer   string `json:"answer" validate:"required"`
}

type GetQuestionPage struct {
	ID             int64  `json:"id"`
	LessonID       int64  `json:"lesson_id"`
	CreatedBy      string `json:"created_by"`
	LastModifiedBy string `json:"last_modified_by"`
	CreatedAt      string `json:"created_at"`
	Modified       string `json:"modified"`
	ContentType    string `json:"content_type"`

	QuestionType string `json:"question_type"`

	Question string `json:"question"`
	OptionA  string `json:"option_a"`
	OptionB  string `json:"option_b"`
	OptionC  string `json:"option_c"`
	OptionD  string `json:"option_d"`
	OptionE  string `json:"option_e"`
	Answer   string `json:"answer"`
}

type UpdateQuestionPage struct {
	ID             int64  `json:"id" validate:"required"`
	ChannelID      int64  `json:"channel_id" validate:"required"`
	PlanID         int64  `json:"plan_id" validate:"required"`
	LessonID       int64  `json:"lesson_id" validate:"required"`
	LastModifiedBy string `json:"last_modified_by" validate:"required"`

	Question string `json:"question,omitempty"`
	OptionA  string `json:"option_a,omitempty"`
	OptionB  string `json:"option_b,omitempty"`
	OptionC  string `json:"option_c,omitempty"`
	OptionD  string `json:"option_d,omitempty"`
	OptionE  string `json:"option_e,omitempty"`
	Answer   string `json:"answer,omitempty"`
}
