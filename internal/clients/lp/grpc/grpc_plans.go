package lpgrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	lpmodels "github.com/DimTur/lp_api_gateway/internal/clients/lp/models"
	"github.com/DimTur/lp_api_gateway/internal/services/permissions.go"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrPlanExitsts  = errors.New("plan already exists")
	ErrPlanNotFound = errors.New("plan not found")
)

func (c *Client) CreatePlan(ctx context.Context, plan *lpmodels.CreatePlan) (*lpmodels.CreatePlanResponse, error) {
	const op = "lp.grpc.CreatePlan"

	resp, err := c.api.CreatePlan(ctx, &lpv1.CreatePlanRequest{
		Name:        plan.Name,
		Description: plan.Description,
		CreatedBy:   plan.CreatedBy,
		ChannelId:   plan.ChannelID,
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

	return &lpmodels.CreatePlanResponse{
		ID:      resp.Id,
		Success: true,
	}, nil
}

func (c *Client) GetPlan(ctx context.Context, plan *lpmodels.GetPlan) (*lpmodels.GetPlanResponse, error) {
	const op = "lp.grpc.GetPlan"

	resp, err := c.api.GetPlan(ctx, &lpv1.GetPlanRequest{
		PlanId:    plan.PlanID,
		ChannelId: plan.ChannelID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("plan not found", slog.String("err", err.Error()))
			return &lpmodels.GetPlanResponse{}, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.GetPlanResponse{}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.GetPlanResponse{
		Id:             resp.Plan.Id,
		Name:           resp.Plan.Name,
		Description:    resp.Plan.Description,
		CreatedBy:      resp.Plan.CreatedBy,
		LastModifiedBy: resp.Plan.LastModifiedBy,
		IsPublished:    resp.Plan.IsPublished,
		Public:         resp.Plan.Public,
		CreatedAt:      resp.Plan.CreatedAt,
		Modified:       resp.Plan.Modified,
	}, nil
}

func (c *Client) GetPlans(ctx context.Context, inputParam *lpmodels.GetPlans) ([]lpmodels.GetPlanResponse, error) {
	const op = "lp.grpc.GetPlans"

	resp, err := c.api.GetPlans(ctx, &lpv1.GetPlansRequest{
		UserId:    inputParam.UserID,
		ChannelId: inputParam.ChannelID,
		Limit:     inputParam.Limit,
		Offset:    inputParam.Offset,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("plans not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	var planResp []lpmodels.GetPlanResponse
	for _, plan := range resp.Plans {
		planResp = append(planResp, lpmodels.GetPlanResponse{
			Id:             plan.Id,
			Name:           plan.Name,
			Description:    plan.Description,
			CreatedBy:      plan.CreatedBy,
			LastModifiedBy: plan.LastModifiedBy,
			IsPublished:    plan.IsPublished,
			Public:         plan.Public,
			CreatedAt:      plan.CreatedAt,
			Modified:       plan.Modified,
		})
	}

	return planResp, nil
}

func (c *Client) UpdatePlan(ctx context.Context, updPlan *lpmodels.UpdatePlan) (*lpmodels.UpdatePlanResponse, error) {
	const op = "lp.grpc.UpdatePlan"

	resp, err := c.api.UpdatePlan(ctx, &lpv1.UpdatePlanRequest{
		ChannelId:      updPlan.ChannelID,
		PlanId:         updPlan.PlanID,
		Name:           updPlan.Name,
		Description:    updPlan.Description,
		LastModifiedBy: updPlan.LastModifiedBy,
		IsPublished:    updPlan.IsPublished,
		Public:         updPlan.Public,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return &lpmodels.UpdatePlanResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.NotFound:
			c.log.Error("plan not found", slog.String("err", err.Error()))
			return &lpmodels.UpdatePlanResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.UpdatePlanResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.UpdatePlanResponse{
		ID:      resp.Id,
		Success: true,
	}, nil
}

func (c *Client) DeletePlan(ctx context.Context, delPlan *lpmodels.DelPlan) (*lpmodels.DelPlanResponse, error) {
	const op = "lp.grpc.DeletePlan"

	resp, err := c.api.DeletePlan(ctx, &lpv1.DeletePlanRequest{
		PlanId:    delPlan.PlanID,
		ChannelId: delPlan.ChannelID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("plan not found", slog.String("err", err.Error()))
			return &lpmodels.DelPlanResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.DelPlanResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.DelPlanResponse{
		Success: resp.Success,
	}, nil
}

func (c *Client) SharePlanWithUser(ctx context.Context, sharePlanWithUser *lpmodels.SharePlan) (*lpmodels.SharingPlanResp, error) {
	const op = "lp.grpc.SharePlanWithUser"

	resp, err := c.api.SharePlanWithUsers(ctx, &lpv1.SharePlanWithUsersRequest{
		ChannelId: sharePlanWithUser.ChannelID,
		PlanId:    sharePlanWithUser.PlanID,
		UsersIds:  sharePlanWithUser.UsersIDs,
		CreatedBy: sharePlanWithUser.UserID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.SharingPlanResp{
		Success: resp.Success,
	}, nil
}

func (c *Client) IsUserShareWithPlan(ctx context.Context, userPlan *permissions.IsUserShareWithPlan) (*lpmodels.IsPlanShareWith, error) {
	const op = "lp.grpc.SharePlanWithUser"

	resp, err := c.api.IsUserShareWithPlan(ctx, &lpv1.IsUserShareWithPlanRequest{
		UserId: userPlan.UserID,
		PlanId: userPlan.PlanID,
	})
	if err != nil {
		switch status.Code(err) {
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.IsPlanShareWith{
				IsShare: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.IsPlanShareWith{
		IsShare: resp.IsShare,
	}, nil
}
