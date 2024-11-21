package pageshandler

type CreateImagePageRequest struct {
	ImageFileUrl string `json:"image_file_url" validate:"required"`
	ImageName    string `json:"image_name" validate:"required"`
}

type CreateVideoPageRequest struct {
	VideoFileUrl string `json:"video_file_url" validate:"required"`
	VideoName    string `json:"video_name" validate:"required"`
}
type CreatePDFPageRequest struct {
	PdfFileUrl string `json:"pdf_file_url" validate:"required"`
	PdfName    string `json:"pdf_name" validate:"required"`
}

type UpdateImagePageRequest struct {
	ImageFileUrl string `json:"image_file_url,omitempty"`
	ImageName    string `json:"image_name,omitempty"`
}

type UpdateVideoPageRequest struct {
	VideoFileUrl string `json:"video_file_url,omitempty"`
	VideoName    string `json:"video_name,omitempty"`
}

type UpdatePDFPageRequest struct {
	PdfFileUrl string `json:"pdf_file_url,omitempty"`
	PdfName    string `json:"pdf_name,omitempty"`
}
