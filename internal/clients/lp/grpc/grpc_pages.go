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
	ErrPageNotFound = errors.New("page not found")
	ErrUnContType   = errors.New("unsupported content type")
)

func (c *Client) CreateImagePage(ctx context.Context, page *lpmodels.CreateImagePage) (*lpmodels.CreatePageResponse, error) {
	const op = "lp.grpc.CreateImagePage"

	resp, err := c.api.CreateImagePage(ctx, &lpv1.CreateImagePageRequest{
		Base: &lpv1.CreateBasePage{
			LessonId:  page.LessonID,
			CreatedBy: page.CreatedBy,
		},
		ImageFileUrl: page.ImageFileUrl,
		ImageName:    page.ImageName,
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

func (c *Client) CreateVideoPage(ctx context.Context, page *lpmodels.CreateVideoPage) (*lpmodels.CreatePageResponse, error) {
	const op = "lp.grpc.CreateVideoPage"

	resp, err := c.api.CreateVideoPage(ctx, &lpv1.CreateVideoPageRequest{
		Base: &lpv1.CreateBasePage{
			LessonId:  page.LessonID,
			CreatedBy: page.CreatedBy,
		},
		VideoFileUrl: page.VideoFileUrl,
		VideoName:    page.VideoName,
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

func (c *Client) CreatePDFPage(ctx context.Context, page *lpmodels.CreatePDFPage) (*lpmodels.CreatePageResponse, error) {
	const op = "lp.grpc.CreateVideoPage"

	resp, err := c.api.CreatePDFPage(ctx, &lpv1.CreatePDFPageRequest{
		Base: &lpv1.CreateBasePage{
			LessonId:  page.LessonID,
			CreatedBy: page.CreatedBy,
		},
		PdfFileUrl: page.PdfFileUrl,
		PdfName:    page.PdfName,
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

func (c *Client) GetImagePage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.ImagePage, error) {
	const op = "lp.grpc.GetImagePage"

	resp, err := c.api.GetImagePage(ctx, &lpv1.GetImagePageRequest{
		PageId:   page.PageID,
		LessonId: page.LessonID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("image page not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.ImagePage{
		BasePage: lpmodels.BasePage{
			ID:             resp.Base.Id,
			LessonID:       resp.Base.LessonId,
			CreatedBy:      resp.Base.CreatedBy,
			LastModifiedBy: resp.Base.LastModifiedBy,
			CreatedAt:      resp.Base.CreatedAt,
			Modified:       resp.Base.Modified,
			ContentType:    resp.Base.ContentType.String(),
		},
		ImageFileUrl: resp.ImageFileUrl,
		ImageName:    resp.ImageName,
	}, nil
}

func (c *Client) GetVideoPage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.VideoPage, error) {
	const op = "lp.grpc.GetVideoPage"

	resp, err := c.api.GetVideoPage(ctx, &lpv1.GetVideoPageRequest{
		PageId:   page.PageID,
		LessonId: page.LessonID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("video page not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.VideoPage{
		BasePage: lpmodels.BasePage{
			ID:             resp.Base.Id,
			LessonID:       resp.Base.LessonId,
			CreatedBy:      resp.Base.CreatedBy,
			LastModifiedBy: resp.Base.LastModifiedBy,
			CreatedAt:      resp.Base.CreatedAt,
			Modified:       resp.Base.Modified,
			ContentType:    resp.Base.ContentType.String(),
		},
		VideoFileUrl: resp.VideoFileUrl,
		VideoName:    resp.VideoName,
	}, nil
}

func (c *Client) GetPDFPage(ctx context.Context, page *lpmodels.GetPage) (*lpmodels.PDFPage, error) {
	const op = "lp.grpc.GetPDFPage"

	resp, err := c.api.GetPDFPage(ctx, &lpv1.GetPDFPageRequest{
		PageId:   page.PageID,
		LessonId: page.LessonID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("pdf page not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.PDFPage{
		BasePage: lpmodels.BasePage{
			ID:             resp.Base.Id,
			LessonID:       resp.Base.LessonId,
			CreatedBy:      resp.Base.CreatedBy,
			LastModifiedBy: resp.Base.LastModifiedBy,
			CreatedAt:      resp.Base.CreatedAt,
			Modified:       resp.Base.Modified,
			ContentType:    resp.Base.ContentType.String(),
		},
		PdfFileUrl: resp.PdfFileUrl,
		PdfName:    resp.PdfName,
	}, nil
}

func (c *Client) GetPages(ctx context.Context, inputParams *lpmodels.GetPages) ([]lpmodels.BasePage, error) {
	const op = "lp.grpc.GetPages"

	resp, err := c.api.GetPages(ctx, &lpv1.GetPagesRequest{
		LessonId: inputParams.LessonID,
		Limit:    inputParams.Limit,
		Offset:   inputParams.Offset,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("pages not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrPageNotFound)
		case codes.InvalidArgument:
			c.log.Error("bad request", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	var pagesResp []lpmodels.BasePage
	for _, page := range resp.Pages {
		pagesResp = append(pagesResp, lpmodels.BasePage{
			ID:             page.Id,
			LessonID:       page.LessonId,
			CreatedBy:      page.CreatedBy,
			LastModifiedBy: page.LastModifiedBy,
			CreatedAt:      page.CreatedAt,
			Modified:       page.Modified,
			ContentType:    page.ContentType.String(),
		})
	}

	return pagesResp, nil
}

func (c *Client) UpdateImagePage(ctx context.Context, updIPage *lpmodels.UpdateImagePage) (*lpmodels.UpdatePageResponse, error) {
	const op = "lp.grpc.UpdateImagePage"

	resp, err := c.api.UpdateImagePage(ctx, &lpv1.UpdateImagePageRequest{
		Base: &lpv1.UpdateBasePage{
			Id:             updIPage.UpdateBasePage.ID,
			LastModifiedBy: updIPage.LastModifiedBy,
		},
		ImageFileUrl: updIPage.ImageFileUrl,
		ImageName:    updIPage.ImageName,
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
			c.log.Error("image page not found", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
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

func (c *Client) UpdateVideoPage(ctx context.Context, updIPage *lpmodels.UpdateVideoPage) (*lpmodels.UpdatePageResponse, error) {
	const op = "lp.grpc.UpdateVideoPage"

	resp, err := c.api.UpdateVideoPage(ctx, &lpv1.UpdateVideoPageRequest{
		Base: &lpv1.UpdateBasePage{
			Id:             updIPage.UpdateBasePage.ID,
			LastModifiedBy: updIPage.LastModifiedBy,
		},
		VideoFileUrl: updIPage.VideoFileUrl,
		VideoName:    updIPage.VideoName,
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
			c.log.Error("video page not found", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
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

func (c *Client) UpdatePDFPage(ctx context.Context, updIPage *lpmodels.UpdatePDFPage) (*lpmodels.UpdatePageResponse, error) {
	const op = "lp.grpc.UpdatePDFPage"

	resp, err := c.api.UpdatePDFPage(ctx, &lpv1.UpdatePDFPageRequest{
		Base: &lpv1.UpdateBasePage{
			Id:             updIPage.UpdateBasePage.ID,
			LastModifiedBy: updIPage.LastModifiedBy,
		},
		PdfFileUrl: updIPage.PdfFileUrl,
		PdfName:    updIPage.PdfName,
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
			c.log.Error("pdf page not found", slog.String("err", err.Error()))
			return &lpmodels.UpdatePageResponse{
				ID:      0,
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrPageNotFound)
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

func (c *Client) DeletePage(ctx context.Context, delPage *lpmodels.DeletePage) (*lpmodels.DeletePageResponse, error) {
	const op = "lp.grpc.DeletePage"

	resp, err := c.api.DeletePage(ctx, &lpv1.DeletePageRequest{
		PageId:   delPage.PageID,
		LessonId: delPage.LessonID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("page not found", slog.String("err", err.Error()))
			return &lpmodels.DeletePageResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return &lpmodels.DeletePageResponse{
				Success: false,
			}, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &lpmodels.DeletePageResponse{
		Success: resp.Success,
	}, nil
}
