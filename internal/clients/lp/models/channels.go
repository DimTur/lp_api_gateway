package lpmodels

import (
	"time"
)

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
	ChannelID int64 `json:"channel_id" validate:"required"`
}

type GetChannelResponse struct {
	Id             int64
	Name           string
	Description    string
	CreatedBy      string
	LastModifiedBy string
	CreatedAt      time.Time
	Modified       time.Time
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
	CreatedAt      time.Time
	Modified       time.Time
}
