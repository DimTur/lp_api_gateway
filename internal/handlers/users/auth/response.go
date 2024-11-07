package authhandler

import "github.com/DimTur/lp_api_gateway/internal/lib/api/response"

type SingUpResponse struct {
	response.Response
	Success bool
}

type SingInResponse struct {
	response.Response
	AccsessToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
