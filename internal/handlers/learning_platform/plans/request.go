package planshandler

type CreatePlanRequest struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description,omitempty"`
	LearningGroupId string `json:"learning_group_id" validate:"required"`
}

type UpdatePlanRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsPublished *bool   `json:"is_published,omitempty"`
	Public      *bool   `json:"public,omitempty"`
}

type SharePlanRequest struct {
	UserIDs []string `json:"user_ids" validate:"required"`
}
