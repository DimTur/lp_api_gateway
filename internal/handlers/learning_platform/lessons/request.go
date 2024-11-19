package lessonshandler

type CreateLessonRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateLessonRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
