package channelshandler

type CreateChannelRequest struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description,omitempty"`
	LearningGroupId string `json:"learning_group_id" validate:"required"`
}

type GetChannelsRequest struct {
	Limit  int64 `json:"limit,omitempty" validate:"min=1"`
	Offset int64 `json:"offset,omitempty" validate:"min=0"`
}

type UpdateChannelRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type ShareChannelRequest struct {
	LGroupIDs []string `json:"lgroup_ids" validate:"required"`
}
