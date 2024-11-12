package lpmodels

type CreateChannel struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	CreatedBy   string `json:"created_by" validate:"required"`
}

type CreateChannelResponse struct {
	ID      int64
	Success bool
}

type GetChannel struct {
	ChannelID        int64    `json:"channel_id" validate:"required"`
	LearningGroupIds []string `json:"learning_group_ids" validate:"required"`
}

type GetChannelResponse struct {
	Id             int64
	Name           string
	Description    string
	CreatedBy      string
	LastModifiedBy string
	CreatedAt      string
	Modified       string
	Plans          []*Plan
}

type Plan struct {
	Id             int64
	Name           string
	Description    string
	CreatedBy      string
	LastModifiedBy string
	IsPublished    bool
	Public         bool
	CreatedAt      string
	Modified       string
}

type GetChannels struct {
	LearningGroupIds []string `json:"learning_group_ids" validate:"required"`
	Limit            int64    `json:"limit,omitempty" validate:"min=1"`
	Offset           int64    `json:"offset,omitempty" validate:"min=0"`
}

type Channel struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      string
	LastModifiedBy string
	CreatedAt      string
	Modified       string
}

type UpdateChannel struct {
	UserID       string   `json:"user_id" validate:"required"`
	AdminInLgIds []string `json:"admin_in_lg_ch_ids" validate:"required"`
	ChannelID    int64    `json:"id" validate:"required"`
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
}

type UpdateChannelResponse struct {
	ID      int64
	Success bool
}

type DelChByID struct {
	ChannelID    int64    `json:"id" validate:"required"`
	AdminInLgIds []string `json:"admin_in_lg_ch_ids" validate:"required"`
}

type DelChByIDResp struct {
	Success bool
}

type SharingChannel struct {
	ChannelID int64    `json:"channel_id" validate:"required"`
	LGroupIDs []string `json:"lgroup_ids" validate:"required"`
	CreatedBy string   `json:"created_by" validate:"required"`
}

type SharingChannelResp struct {
	Success bool
}
