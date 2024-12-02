package lpgrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrLessonAttemtNotFound       = errors.New("lesson attempt not found")
	ErrQuestionPageAttemtNotFound = errors.New("question page attempt not found")
	ErrAnswerNotFound             = errors.New("page answer not found")
	ErrPermissionsDenied          = errors.New("permissions denied")
)

func (c *Client) TryLesson(ctx context.Context, lesson *lpmodels.TryLesson) (*lpmodels.TryLessonResp, error) {
	const op = "lp.grpc.TryLesson"

	resp, err := c.api.TryLesson(ctx, &lpv1.TryLessonRequest{
		UserId:    lesson.UserID,
		LessonId:  lesson.LessonID,
		PlanId:    lesson.PlanID,
		ChannelId: lesson.ChannelID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("question page attempt not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrQuestionPageAttemtNotFound)
		case codes.InvalidArgument:
			c.log.Error("invalid arguments", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	var attemtResp lpmodels.TryLessonResp
	for _, attempt := range resp.QuestionPageAttempts {
		attemtResp.QuestionPageAttempts = append(attemtResp.QuestionPageAttempts, lpmodels.QuestionPageAttempt{
			ID:              attempt.Id,
			PageID:          attempt.PageId,
			LessonAttemptID: attempt.LessonAttemptId,
			IsCorrect:       attempt.IsCorrect,
			UserAnswer:      attempt.UserAnswer.Enum().String(),
		})
	}

	return &attemtResp, nil
}

func (c *Client) UpdatePageAttempt(ctx context.Context, attempt *lpmodels.UpdatePageAttempt) (*lpmodels.UpdatePageAttemptResp, error) {
	const op = "lp.grpc.UpdatePageAttempt"

	resp, err := c.api.UpdatePageAttempt(ctx, &lpv1.UpdatePageAttemptRequest{
		QuestionAttemptId: attempt.QPAttemptID,
		PageId:            attempt.PageID,
		LessonAttemptId:   attempt.LessonAttemptID,
		UserAnswer:        lpv1.Answer(attempt.PageID),
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.NotFound:
			c.log.Error("question page attempt answer not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrAnswerNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.UpdatePageAttemptResp{
		Success: resp.Success,
	}, nil
}

func (c *Client) CompleteLesson(ctx context.Context, lesson *lpmodels.CompleteLesson) (*lpmodels.CompleteLessonResp, error) {
	const op = "lp.grpc.CompleteLesson"

	resp, err := c.api.CompleteLesson(ctx, &lpv1.CompleteLessonRequest{
		UserId:          lesson.UserID,
		LessonAttemptId: lesson.LessonAttemptID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("lesson page attempt not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrQuestionPageAttemtNotFound)
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.CompleteLessonResp{
		ID:              resp.LessonAttemptId,
		IsSuccessful:    resp.IsSuccessfull,
		PercentageScore: resp.PercentageScore,
	}, nil
}

func (c *Client) GetLessonAttempts(ctx context.Context, inputParams *lpmodels.GetLessonAttempts) (*lpmodels.GetLessonAttemptsResp, error) {
	const op = "lp.grpc.GetLessonAttempts"

	resp, err := c.api.GetLessonAttempts(ctx, &lpv1.GetLessonAttemptsRequest{
		UserId:   inputParams.UserID,
		LessonId: inputParams.LessonID,
		Limit:    inputParams.Limit,
		Offset:   inputParams.Offset,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("lesson attempt not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrLessonAttemtNotFound)
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	var attempResp lpmodels.GetLessonAttemptsResp
	for _, attempt := range resp.LessonAttempts {
		attempResp.LessonAttempts = append(attempResp.LessonAttempts, lpmodels.LessonAttempt{
			ID:              attempt.Id,
			UserID:          attempt.UserId,
			LessonID:        attempt.LessonId,
			PlanID:          attempt.PlanId,
			ChannelID:       attempt.ChannelId,
			StartTime:       attempt.StartTime,
			EndTime:         attempt.EndTime,
			IsComplete:      attempt.IsComplete,
			IsSuccessful:    attempt.IsSuccessful,
			PercentageScore: attempt.PercentageScore,
		})
	}

	return &attempResp, nil
}

func (c *Client) CheckLessonAttemptPermissions(ctx context.Context, userAtt *lpmodels.LessonAttemptPermissions) (bool, error) {
	const op = "lp.grpc.GetLessonAttempts"

	resp, err := c.api.CheckPermissionForUser(ctx, &lpv1.CheckPermissionForUserRequest{
		UserId:          userAtt.UserID,
		LessonAttemptId: userAtt.LessonAttemptID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.PermissionDenied:
			c.log.Error("permission denied", slog.String("err", err.Error()))
			return false, fmt.Errorf("%s: %w", op, ErrPermissionsDenied)
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return false, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return resp.Success, nil
}
