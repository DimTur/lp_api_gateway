package authhandler

import "github.com/DimTur/lp_api_gateway/internal/lib/api/response"

type SingUpResponse struct {
	response.Response
	Success bool
}

type SingInResponse struct {
	response.Response
	AccsessToken string
	RefreshToken string
}

type SingInByTgResponse struct {
	response.Response
	Success bool
	Info    string
}

type CheckOTPAndLogInResponse struct {
	response.Response
	AccsessToken string
	RefreshToken string
}

type UpdateUserInfoResponse struct {
	response.Response
	Success bool
}
