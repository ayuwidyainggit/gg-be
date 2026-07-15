package adapter

import (
	"fmt"
	"inventory/model"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type ObsAdapter interface {
	UploadFile(req *model.Upload) (fullUrl string, err error)
	UploadBytes(req *model.UploadBytes) (fullUrl string, err error)
}
type ObsAdapterImpl struct {
	Obsc        *obs.ObsClient
	FileBaseUrl string
	Bucket      string
}

func InitObsAdapter(AK, SK, ENDPOINT, BUCKET string) (*ObsAdapterImpl, error) {
	endpoint := fmt.Sprintf("https://%v", ENDPOINT)
	fileBaseUrl := fmt.Sprintf("https://%v.%v", BUCKET, ENDPOINT)

	obsClient, err := obs.New(AK, SK, endpoint /*, obs.WithSecurityToken(securityToken)*/)
	if err != nil {
		// Use the struct to access OBS.
		fmt.Printf("Create obsClient error, errMsg: %s", err.Error())
		// Close obsClient.
		obsClient.Close()
		return nil, err
	}

	return &ObsAdapterImpl{Obsc: obsClient, FileBaseUrl: fileBaseUrl, Bucket: BUCKET}, nil
}

func (o *ObsAdapterImpl) UploadFile(req *model.Upload) (fullUrl string, err error) {
	key, err := req.GenerateFileName() // generate file name
	if err != nil {
		return
	}

	byteReader, err := req.FileConvertToByteReader() // convert form to byte reader
	if err != nil {
		return
	}

	input := &obs.PutObjectInput{}
	// Specify a bucket name.
	input.Bucket = o.Bucket
	// Specify the object (example/objectname as an example) to upload.
	input.Key = key
	input.ACL = obs.AclType(obs.AclPublicRead)
	input.ContentType = req.GetFileContentType()
	input.Body = byteReader
	// Upload you local file using streaming.
	_, err = o.Obsc.PutObject(input)
	if err == nil {
		fullUrl = fmt.Sprintf("%v/%v", o.FileBaseUrl, key)
		return
	}
	if obsError, ok := err.(obs.ObsError); ok {
		return fullUrl, obsError
	}
	return
}

// UploadBytes uploads raw byte data to OBS using provided metadata.
func (o *ObsAdapterImpl) UploadBytes(req *model.UploadBytes) (fullUrl string, err error) {
	key, err := req.GenerateFileName()
	if err != nil {
		return
	}

	byteReader, err := req.FileConvertToByteReader()
	if err != nil {
		return
	}

	input := &obs.PutObjectInput{}
	input.Bucket = o.Bucket
	input.Key = key
	input.ACL = obs.AclType(obs.AclPublicRead)
	input.ContentType = req.ContentType
	input.Body = byteReader
	_, err = o.Obsc.PutObject(input)
	if err == nil {
		fullUrl = fmt.Sprintf("%v/%v", o.FileBaseUrl, key)
		return
	}
	if obsError, ok := err.(obs.ObsError); ok {
		return fullUrl, obsError
	}
	return
}
