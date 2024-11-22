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
	ErrInvalidAnswer    = errors.New("invalid answer value")
	ErrQuestionNotFound = errors.New("question not found")
)

func (c *Client) CreateQuestionPage(ctx context.Context, question *lpmodels.CreateQuestionPage) (*lpmodels.CreatePageResponse, error) {
	const op = "lp.grpc.CreateQuestionPage"

	answerEnum, err := toAnswerEnum(question.Answer)
	if err != nil {
		c.log.Error("invalid answer value", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidAnswer)
	}

	resp, err := c.api.CreateQuestionPage(ctx, &lpv1.CreateQuestionPageRequest{
		LessonId:  question.LessonID,
		CreatedBy: question.CreatedBy,
		Question:  question.Question,
		OptionA:   question.OptionA,
		OptionB:   question.OptionB,
		OptionC:   &question.OptionC,
		OptionD:   &question.OptionD,
		OptionE:   &question.OptionE,
		Answer:    answerEnum,
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

	return &lpmodels.CreatePageResponse{
		ID:      resp.Id,
		Success: true,
	}, nil
}

func (c *Client) GetQuestionPage(ctx context.Context, question *lpmodels.GetPage) (*lpmodels.GetQuestionPage, error) {
	const op = "lp.grpc.GetQuestionPage"

	resp, err := c.api.GetQuestionPage(ctx, &lpv1.GetQuestionPageRequest{
		PageId:   question.PageID,
		LessonId: question.LessonID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("invalid arguments", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.NotFound:
			c.log.Error("question page not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrQuestionNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.GetQuestionPage{
		ID:             resp.QuestionPage.Id,
		LessonID:       resp.QuestionPage.LessonId,
		CreatedBy:      resp.QuestionPage.CreatedBy,
		LastModifiedBy: resp.QuestionPage.LastModifiedBy,
		CreatedAt:      resp.QuestionPage.CreatedAt,
		Modified:       resp.QuestionPage.Modified,
		ContentType:    resp.QuestionPage.ContentType.String(),
		QuestionType:   resp.QuestionPage.QuestionType.String(),
		Question:       resp.QuestionPage.Question,
		OptionA:        resp.QuestionPage.OptionA,
		OptionB:        resp.QuestionPage.OptionB,
		OptionC:        resp.QuestionPage.OptionC,
		OptionD:        resp.QuestionPage.OptionD,
		OptionE:        resp.QuestionPage.OptionE,
		Answer:         resp.QuestionPage.Answer,
	}, nil

}

func (c *Client) UpdateQuestionPage(ctx context.Context, updQust *lpmodels.UpdateQuestionPage) (*lpmodels.UpdatePageResponse, error) {
	const op = "lp.grpc.UpdateQuestionPage"

	answerEnum, err := toAnswerEnum(updQust.Answer)
	if err != nil {
		c.log.Error("invalid answer value", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidAnswer)
	}

	resp, err := c.api.UpdateQuestionPage(ctx, &lpv1.UpdateQuestionPageRequest{
		Id:             updQust.ID,
		LastModifiedBy: updQust.LastModifiedBy,
		Question:       &updQust.Question,
		OptionA:        &updQust.OptionA,
		OptionB:        &updQust.OptionB,
		OptionC:        &updQust.OptionC,
		OptionD:        &updQust.OptionD,
		OptionE:        &updQust.OptionE,
		Answer:         &answerEnum,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.NotFound:
			c.log.Error("question page not found", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrQuestionNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.UpdatePageResponse{
		ID:      resp.Id,
		Success: true,
	}, nil
}

func toAnswerEnum(answer string) (lpv1.Answer, error) {
	switch answer {
	case "OPTION_A":
		return lpv1.Answer_OPTION_A, nil
	case "OPTION_B":
		return lpv1.Answer_OPTION_B, nil
	case "OPTION_C":
		return lpv1.Answer_OPTION_C, nil
	case "OPTION_D":
		return lpv1.Answer_OPTION_D, nil
	case "OPTION_E":
		return lpv1.Answer_OPTION_E, nil
	default:
		return lpv1.Answer_ANSWER_UNSPECIFIED, fmt.Errorf("invalid answer value: %s", answer)
	}
}
