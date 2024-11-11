package learninggrouphandler

import (
	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
)

type CreateLGroupResponse struct {
	response.Response
	Success bool
}

type GetLgByIDResponse struct {
	response.Response
	LearningGroup *ssomodels.GetLgByIDResp
}

type UpdateLGroupResponse struct {
	response.Response
	Success bool
}

type DeleteLGroupResponse struct {
	response.Response
	Success bool
}

type GetLearningGroupsResponse struct {
	response.Response
	LearningGroups *ssomodels.GetLGroupsResp
}
