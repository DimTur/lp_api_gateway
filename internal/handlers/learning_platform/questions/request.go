package questionshandler

type CreateQuestionPageRequest struct {
	Question string `json:"question" validate:"required"`
	OptionA  string `json:"option_a" validate:"required"`
	OptionB  string `json:"option_b" validate:"required"`
	OptionC  string `json:"option_c,omitempty"`
	OptionD  string `json:"option_d,omitempty"`
	OptionE  string `json:"option_e,omitempty"`
	Answer   string `json:"answer" validate:"required"`
}

type UpdateQuestinPageRequest struct {
	Question string `json:"question,omitempty"`
	OptionA  string `json:"option_a,omitempty"`
	OptionB  string `json:"option_b,omitempty"`
	OptionC  string `json:"option_c,omitempty"`
	OptionD  string `json:"option_d,omitempty"`
	OptionE  string `json:"option_e,omitempty"`
	Answer   string `json:"answer,omitempty"`
}
