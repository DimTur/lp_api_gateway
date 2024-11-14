package permissions

type CheckPerm struct {
	UserID    string `json:"user_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
}
