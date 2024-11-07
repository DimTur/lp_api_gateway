package channelshandler

import (
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
)

type CreateChannelResponce struct {
	response.Response
	ChannelID int64 `json:"channel_id,omitempty"`
}

type GetChannelResponce struct {
	response.Response
	Channel lpmodels.GetChannelResponse
}
