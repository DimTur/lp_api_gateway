package planshandler

import (
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
)

type CreatePlanResponse struct {
	response.Response
	PlanID int64 `json:"plan_id,omitempty"`
}

type GetPlanResponse struct {
	response.Response
	Plan lpmodels.GetPlanResponse
}

type GetPlansResponse struct {
	response.Response
	Plans []lpmodels.GetPlanResponse
}

type UpdatePlanResponse struct {
	response.Response
	UpdatePlanResponse lpmodels.UpdatePlanResponse
}

type DeletePlanResponse struct {
	response.Response
	Success bool
}

type SharePlanResponse struct {
	response.Response
	Success bool
}
