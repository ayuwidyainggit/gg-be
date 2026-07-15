package service

import (
	"sales/adapter"
	"sales/entity"
	"sales/model"
	"sales/pkg/config/env"
)

type FilesService interface {
	Upload(request entity.UploadRequest) (resp entity.UploadResponse, err error)
}

type FilesServiceImpl struct {
	Config env.ConfigEnv
	// MCustomerRepository repository.MCustomerRepository,
	ObsAdapter adapter.ObsAdapter
}

func NewFilesService(
	config env.ConfigEnv,
	obsAdapter adapter.ObsAdapter,
) *FilesServiceImpl {
	return &FilesServiceImpl{
		Config:     config,
		ObsAdapter: obsAdapter,
	}
}

func (service *FilesServiceImpl) Upload(request entity.UploadRequest) (resp entity.UploadResponse, err error) {
	uploadModel := &model.Upload{
		Folder: *request.Folder,
		File:   request.File,
	}

	fileUrl, err := service.ObsAdapter.UploadFile(uploadModel)
	if err != nil {
		return resp, err
	}
	resp.Url = fileUrl

	return
}
