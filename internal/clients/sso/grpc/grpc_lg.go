package ssogrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	ssov1 "github.com/DimTur/lp_protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidGroupID   = errors.New("invalid group id")
	ErrGroupExists      = errors.New("group already exists")
	ErrGroupNotFound    = errors.New("group not found")
	ErrPermissionDenied = errors.New("you don't have permissions")
)

func (c *Client) CreateLearningGroup(ctx context.Context, newLGroup *ssomodels.CreateLearningGroup) (*ssomodels.CreateLearningGroupResp, error) {
	const op = "sso.grpc_lg.CreateLearningGroup"

	resp, err := c.api.CreateLearningGroup(ctx, &ssov1.CreateLearningGroupRequest{
		Name:        newLGroup.Name,
		CreatedBy:   newLGroup.CreatedBy,
		ModifiedBy:  newLGroup.ModifiedBy,
		GroupAdmins: newLGroup.GroupAdmins,
		Learners:    newLGroup.Learners,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.log.Error("invalid arguments", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		case codes.AlreadyExists:
			c.log.Error("group alredy exists", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupExists)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.CreateLearningGroupResp{
		Success: resp.Success,
	}, nil
}

func (c *Client) GetLearningGroupByID(ctx context.Context, lgID *ssomodels.GetLgByID) (*ssomodels.GetLgByIDResp, error) {
	const op = "sso.grpc_lg.GetLearningGroupByID"

	resp, err := c.api.GetLearningGroupByID(ctx, &ssov1.GetLearningGroupByIDRequest{
		UserId:          lgID.UserID,
		LearningGroupId: lgID.LgId,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.PermissionDenied:
			c.log.Error("permissions denied", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
		case codes.NotFound:
			c.log.Error("group not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	lgResp := &ssomodels.GetLgByIDResp{
		Id:          resp.Id,
		Name:        resp.Name,
		CreatedBy:   resp.CreatedBy,
		ModifiedBy:  resp.ModifiedBy,
		GroupAdmins: make([]*ssomodels.GroupAdmins, len(resp.GroupAdmins)),
		Learners:    make([]*ssomodels.Learner, len(resp.Learners)),
	}

	for i, ga := range resp.GroupAdmins {
		lgResp.GroupAdmins[i] = &ssomodels.GroupAdmins{
			Id:    ga.Id,
			Email: ga.Email,
			Name:  ga.Name,
		}
	}

	for i, l := range resp.Learners {
		lgResp.Learners[i] = &ssomodels.Learner{
			Id:    l.Id,
			Email: l.Email,
			Name:  l.Name,
		}
	}

	return lgResp, nil
}

func (c *Client) UpdateLearningGroup(ctx context.Context, updFields *ssomodels.UpdateLearningGroup) (*ssomodels.UpdateLearningGroupResp, error) {
	const op = "sso.grpc_lg.UpdateLearningGroup"

	resp, err := c.api.UpdateLearningGroup(ctx, &ssov1.UpdateLearningGroupRequest{
		UserId:          updFields.UserID,
		LearningGroupId: updFields.LgId,
		Name:            updFields.Name,
		ModifiedBy:      updFields.ModifiedBy,
		GroupAdmins:     updFields.GroupAdmins,
		Learners:        updFields.Learners,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.PermissionDenied:
			c.log.Error("permissions denied", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
		case codes.NotFound:
			c.log.Error("group not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		case codes.InvalidArgument:
			c.log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.UpdateLearningGroupResp{
		Success: resp.Success,
	}, nil
}

func (c *Client) DeleteLearningGroup(ctx context.Context, lgID *ssomodels.DelLgByID) (*ssomodels.DelLgByIDResp, error) {
	const op = "sso.grpc_lg.DeleteLearningGroup"

	resp, err := c.api.DeleteLearningGroup(ctx, &ssov1.DeleteLearningGroupRequest{
		UserId:          lgID.UserID,
		LearningGroupId: lgID.LgID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.PermissionDenied:
			c.log.Error("permissions denied", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionDenied)
		default:
			c.log.Error("invalid credentials", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}

	return &ssomodels.DelLgByIDResp{
		Success: resp.Success,
	}, nil
}

func (c *Client) GetLearningGroups(ctx context.Context, uID *ssomodels.GetLGroups) (*ssomodels.GetLGroupsResp, error) {
	const op = "sso.grpc_lg.DeleteLearningGroup"

	lGroups, err := c.api.GetLearningGroups(ctx, &ssov1.GetLearningGroupsRequest{
		UserId: uID.UserID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("groups not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	resp := &ssomodels.GetLGroupsResp{
		LearningGroups: make([]*ssomodels.LearningGroup, len(lGroups.LearningGroups)),
	}
	for i, g := range lGroups.LearningGroups {
		resp.LearningGroups[i] = &ssomodels.LearningGroup{
			Id:         g.Id,
			Name:       g.Name,
			CreatedBy:  g.CreatedBy,
			ModifiedBy: g.ModifiedBy,
			Created:    g.Created,
			Updated:    g.Updated,
		}
	}

	return resp, nil
}
