package questionshandler

import (
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
)

type CreatePageResponse struct {
	response.Response
	PageID int64
}

type GetQuestionPageResponse struct {
	response.Response
	QuestionPage lpmodels.GetQuestionPage
}

type UpdatePageResponse struct {
	response.Response
	UpdatePageResponse lpmodels.UpdatePageResponse
}
