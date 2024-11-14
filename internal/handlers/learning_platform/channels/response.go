package channelshandler

import (
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
)

type CreateChannelResponse struct {
	response.Response
	ChannelID int64 `json:"channel_id,omitempty"`
}

type GetChannelResponse struct {
	response.Response
	Channel lpmodels.GetChannelResponse
}

type GetChannelsResponse struct {
	response.Response
	Channels []lpmodels.Channel
}

type UpdateChannelResponse struct {
	response.Response
	UpdateChannelResponse lpmodels.UpdateChannelResponse
}

type DeleteChannelResponse struct {
	response.Response
	Success bool
}

type ShareChannelResponse struct {
	response.Response
	Success bool
}
