package model

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"time"

	"github.com/rs/xid"
)

type BucketList struct {
	Name      string
	CreatedAt time.Time
	Location  string
}

type Upload struct {
	Folder string
	File   *multipart.FileHeader
}

// UploadBytes represents binary data to be uploaded without multipart header.
type UploadBytes struct {
	Folder      string
	FileName    string
	Data        []byte
	ContentType string
}

// GenerateFileName for UploadBytes simply returns provided filename if it has an extension.
func (u *UploadBytes) GenerateFileName() (string, error) {
	if u.FileName == "" {
		return "", errors.New("filename is required")
	}
	extension := path.Ext(u.FileName)
	if extension == "" {
		return "", errors.New("file must be has extension")
	}
	return fmt.Sprintf("%v/%v", u.Folder, u.FileName), nil
}

// FileConvertToByteReader for UploadBytes wraps data slice into bytes.Reader.
func (u *UploadBytes) FileConvertToByteReader() (*bytes.Reader, error) {
	return bytes.NewReader(u.Data), nil
}

func (u *Upload) GenerateFileName() (string, error) {
	originFilename := u.File.Filename
	extension := path.Ext(originFilename)
	if extension == "" {
		return "", errors.New("file must be has extension")
	}
	guid := xid.New()
	NewFileName := fmt.Sprintf("%v/%v%v", u.Folder, guid.String(), extension)
	return NewFileName, nil
}

func (u *Upload) GetFileContentType() string {
	fileheader := u.File.Header
	types, ok := fileheader["Content-Type"]
	var contentType string
	if ok {
		// This should be true!
		for _, x := range types {
			contentType = x
			// Most usually you will probably see only one
		}
	}
	return contentType
}

func (u *Upload) FileConvertToByteReader() (*bytes.Reader, error) {
	src, err := u.File.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	size := u.File.Size

	// Read the file into a byte slice
	bs := make([]byte, size)
	_, err = bufio.NewReader(src).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return nil, err
	}
	return bytes.NewReader(bs), nil
}
