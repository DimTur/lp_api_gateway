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
	ErrLessonNotFound = errors.New("lesson not found")
)

func (c *Client) CreateLesson(ctx context.Context, lesson *lpmodels.CreateLesson) (*lpmodels.CreateLessonResponse, error) {
	const op = "lp.grpc.CreateLesson"

	resp, err := c.api.CreateLesson(ctx, &lpv1.CreateLessonRequest{
		Name:        lesson.Name,
		Description: lesson.Description,
		CreatedBy:   lesson.CreatedBy,
		PlanId:      lesson.PlanID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("invalid arguments", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.CreateLessonResponse{
		ID:      resp.Id,
		Success: true,
	}, nil
}

func (c *Client) GetLesson(ctx context.Context, lesson *lpmodels.GetLesson) (*lpmodels.GetLessonResponse, error) {
	const op = "lp.grpc.GetLesson"

	resp, err := c.api.GetLesson(ctx, &lpv1.GetLessonRequest{
		LessonId: lesson.LessonID,
		PlanId:   lesson.PlanID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("lesson not found", slog.String("err", err.Error()))
			return &lpmodels.GetLessonResponse{}, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.GetLessonResponse{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.GetLessonResponse{
		ID:             resp.Lesson.Id,
		Name:           resp.Lesson.Name,
		Description:    resp.Lesson.Description,
		CreatedBy:      resp.Lesson.CreatedBy,
		LastModifiedBy: resp.Lesson.LastModifiedBy,
		CreatedAt:      resp.Lesson.CreatedAt,
		Modified:       resp.Lesson.Modified,
	}, nil
}

func (c *Client) GetLessons(ctx context.Context, inputParam *lpmodels.GetLessons) ([]lpmodels.GetLessonResponse, error) {
	const op = "lp.grpc.GetLesson"

	resp, err := c.api.GetLessons(ctx, &lpv1.GetLessonsRequest{
		PlanId: inputParam.PlanID,
		Limit:  inputParam.Limit,
		Offset: inputParam.Offset,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("lessons not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	var lessonResp []lpmodels.GetLessonResponse
	for _, lesson := range resp.Lessons {
		lessonResp = append(lessonResp, lpmodels.GetLessonResponse{
			ID:             lesson.Id,
			Name:           lesson.Name,
			Description:    lesson.Description,
			CreatedBy:      lesson.CreatedBy,
			LastModifiedBy: lesson.LastModifiedBy,
			CreatedAt:      lesson.CreatedAt,
			Modified:       lesson.Modified,
		})
	}

	return lessonResp, nil
}

func (c *Client) UpdateLesson(ctx context.Context, updLesson *lpmodels.UpdateLesson) (*lpmodels.UpdateLessonResponse, error) {
	const op = "lp.grpc.UpdateLesson"

	resp, err := c.api.UpdateLesson(ctx, &lpv1.UpdateLessonRequest{
		PlanId:         updLesson.PlanID,
		LessonId:       updLesson.LessonID,
		Name:           updLesson.Name,
		Description:    updLesson.Description,
		LastModifiedBy: updLesson.LastModifiedBy,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdateLessonResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.NotFound:
			c.log.Error("lesson not found", slog.String("err", err.Error()))
			return &lpmodels.UpdateLessonResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.UpdateLessonResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.UpdateLessonResponse{
		ID:      resp.Id,
		Success: true,
	}, nil
}

func (c *Client) DeleteLesson(ctx context.Context, delLess *lpmodels.DeleteLesson) (*lpmodels.DeleteLessonResponse, error) {
	const op = "lp.grpc.DeleteLesson"

	resp, err := c.api.DeleteLesson(ctx, &lpv1.DeleteLessonRequest{
		LessonId: delLess.LessonID,
		PlanId:   delLess.PlanID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("lesson not found", slog.String("err", err.Error()))
			return &lpmodels.DeleteLessonResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.DeleteLessonResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.DeleteLessonResponse{
		Success: resp.Success,
	}, nil
}
