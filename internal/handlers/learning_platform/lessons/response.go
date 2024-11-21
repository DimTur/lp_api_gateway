package lessonshandler

import (
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
)

type CreateLessonResponse struct {
	response.Response
	LessonID int64
}

type GetLessonResponse struct {
	response.Response
	Lesson lpmodels.GetLessonResponse
}

type GetLessonsResponse struct {
	response.Response
	Lessons []lpmodels.GetLessonResponse
}

type UpdateLessonResponse struct {
	response.Response
	UpdateLessonResponse lpmodels.UpdateLessonResponse
}

type DeleteLessonResponse struct {
	response.Response
	Success bool
}
