package permissions

type CheckPerm struct {
	UserID    string `json:"user_id" validate:"required"`
	PlanID    int64  `json:"plna_id,omitempty"`
	ChannelID int64  `json:"channel_id" validate:"required"`
}

type IsUserShareWithPlan struct {
	UserID string `json:"user_id" validate:"required"`
	PlanID int64  `json:"plan_id" validate:"required"`
}

type UserShareWithPlans struct {
	UserID string `json:"user_id" validate:"required"`
}
