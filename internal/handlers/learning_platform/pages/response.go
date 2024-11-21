package pageshandler

import (
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
)

type CreatePageResponse struct {
	response.Response
	PageID int64
}

type GetImagePageResponse struct {
	response.Response
	ImagePage lpmodels.ImagePage
}

type GetVideoPageResponse struct {
	response.Response
	VideoPage lpmodels.VideoPage
}

type GetPDFPageResponse struct {
	response.Response
	PDFPage lpmodels.PDFPage
}

type GetPagesResponse struct {
	response.Response
	Pages []lpmodels.BasePage
}

type UpdatePageResponse struct {
	response.Response
	UpdatePageResponse lpmodels.UpdatePageResponse
}

type DeletePageResponse struct {
	response.Response
	Success bool
}
