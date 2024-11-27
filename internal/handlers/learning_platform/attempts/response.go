package attemptshandler

import (
	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/lib/api/response"
)

type TryLessonResponse struct {
	response.Response
	QuestionPageAttempts []lpmodels.QuestionPageAttempt
}

type UpdatePageAttemptResponse struct {
	response.Response
	Success bool
}

type CompleteLessonResponse struct {
	response.Response
	ID              int64
	IsSuccessful    bool
	PercentageScore int64
}

type LessonAttemptsResponse struct {
	response.Response
	LessonAttempts []lpmodels.LessonAttempt
}
