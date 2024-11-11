package learninggrouphandler

type CreateLearningGroupRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type GetLgByIDRequest struct {
	LgID string `json:"learning_group_id" validate:"required"`
}

type UpdateLearningGroupRequest struct {
	Name        string   `json:"name,omitempty"`
	GroupAdmins []string `json:"group_admins,omitempty"`
	Learners    []string `json:"learners,omitempty"`
}

type DelLgByIDRequest struct {
	LgID string `json:"learning_group_id" validate:"required"`
}
