package entity

import "mime/multipart"

type UploadRequest struct {
	Folder *string `form:"folder" validate:"required"`
	File   *multipart.FileHeader
}
type UploadResponse struct {
	Url string `json:"url"`
}
