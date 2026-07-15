package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"inventory/pkg/structs"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	logrus "github.com/sirupsen/logrus"
)

type HttpClientInfo struct {
	JobID   string
	Method  string
	Url     string
	Auth    string
	Payload interface{}
	Headers map[string]string
	Timeout int
	IsDebug bool
}

func (i *HttpClientInfo) prepare() *resty.Request {
	restyInstance := resty.New()

	if i.Timeout > 0 {
		restyInstance.SetTimeout(time.Duration(i.Timeout * int(time.Second)))
	} else {
		restyInstance.SetTimeout(time.Duration(1 * time.Minute))
	}

	restyInstance.SetDebug(i.IsDebug)

	req := restyInstance.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	if i.Auth != "" {
		req.SetHeader("Authorization", i.Auth)
	}

	//set header based on the request
	req.SetHeaders(i.Headers)

	return req
}

func (i *HttpClientInfo) Dispatch(out interface{}) (*resty.Response, error) {
	var (
		err  error
		resp *resty.Response
	)

	req := i.prepare()
	if i.Method == "" {
		return nil, errors.New("missing http method for dispatch request")
	}
	if i.Auth != "" {
		req.SetHeader("Authorization", i.Auth)
	}

	i.Method = strings.ToUpper(i.Method)

	switch i.Method {
	case http.MethodGet:
		if i.Payload != nil {
			if params, ok := i.Payload.(map[string]interface{}); ok {
				values := url.Values{}
				for k, v := range params {
					switch val := v.(type) {
					case []string:
						for _, item := range val {
							values.Add(k, item)
						}
					case []int:
						for _, item := range val {
							values.Add(k, strconv.Itoa(item))
						}
					case []int64:
						for _, item := range val {
							values.Add(k, strconv.FormatInt(item, 10))
						}
					case []interface{}:
						for _, item := range val {
							values.Add(k, fmt.Sprintf("%v", item))
						}
					default:
						values.Add(k, fmt.Sprintf("%v", v))
					}
				}
				req.SetQueryParamsFromValues(values)
			}
		}
		resp, err = req.Get(i.Url)

	case http.MethodPost:
		resp, err = req.SetBody(i.Payload).Post(i.Url)
	case http.MethodPut:
		resp, err = req.SetBody(i.Payload).Put(i.Url)
	case http.MethodPatch:
		resp, err = req.SetBody(i.Payload).Patch(i.Url)
	case http.MethodDelete:
		if i.Payload != nil {
			req.SetBody(i.Payload)
		}
		resp, err = req.Delete(i.Url)
	default:
		if i.Payload == nil {
			resp, err = req.Execute(i.Method, i.Url)
		} else {
			resp, err = req.SetBody(i.Payload).Execute(i.Method, i.Url)
		}
	}

	if err != nil {
		return nil, err
	}

	// ✅ kalau ada output struct, otomatis decode JSON-nya
	if out != nil {
		if unmarshalErr := json.Unmarshal(resp.Body(), out); unmarshalErr != nil {
			return resp, fmt.Errorf("failed to unmarshal response: %w", unmarshalErr)
		}
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithFields(logrus.Fields{
		"url":     req.URL,
		"method":  i.Method,
		"jobId":   i.JobID,
		"reqBody": structs.StructToJson(i.Payload),
		"resBody": resp.String(),
		"status":  resp.StatusCode(),
		"time":    time.Now(),
		"latency": resp.Time(),
	}).Infoln("Outgoing Request")

	return resp, err
}

// MultipartFileInfo represents file information for multipart upload
type MultipartFileInfo struct {
	FieldName string
	FileName  string
	File      *multipart.FileHeader
}

// DispatchMultipart uploads multipart form data with files
func (i *HttpClientInfo) DispatchMultipart(files []MultipartFileInfo, formData map[string]string, out interface{}) (*resty.Response, error) {
	var (
		err  error
		resp *resty.Response
	)

	restyInstance := resty.New()

	if i.Timeout > 0 {
		restyInstance.SetTimeout(time.Duration(i.Timeout * int(time.Second)))
	} else {
		restyInstance.SetTimeout(time.Duration(1 * time.Minute))
	}

	restyInstance.SetDebug(i.IsDebug)

	req := restyInstance.R().
		SetHeader("Accept", "application/json")

	if i.Auth != "" {
		req.SetHeader("Authorization", i.Auth)
	}

	// Set custom headers
	req.SetHeaders(i.Headers)

	if i.Method == "" {
		return nil, errors.New("missing http method for dispatch request")
	}

	i.Method = strings.ToUpper(i.Method)

	// Set form data
	if formData != nil {
		req.SetFormData(formData)
	}

	// Set files - open all files first, then set them to request
	// If error occurs before request is sent, close all opened files
	var openedFiles []multipart.File
	defer func() {
		// Close all opened files if error occurred before request was sent
		// (resty will handle closing files after request completes successfully)
		if len(openedFiles) > 0 {
			for _, f := range openedFiles {
				if closeErr := f.Close(); closeErr != nil {
					// Log but don't fail - we're already handling an error or resty will close it
				}
			}
		}
	}()

	for _, fileInfo := range files {
		if fileInfo.File != nil {
			// Open file from multipart.FileHeader
			file, openErr := fileInfo.File.Open()
			if openErr != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", fileInfo.FileName, openErr)
			}
			openedFiles = append(openedFiles, file)
			// Note: resty will handle closing the file after request completes
			req.SetFileReader(fileInfo.FieldName, fileInfo.FileName, file)
		}
	}

	// Execute request based on method
	switch i.Method {
	case http.MethodPost:
		resp, err = req.Post(i.Url)
	case http.MethodPut:
		resp, err = req.Put(i.Url)
	case http.MethodPatch:
		resp, err = req.Patch(i.Url)
	default:
		resp, err = req.Execute(i.Method, i.Url)
	}

	// Clear openedFiles after request is sent - resty will handle closing
	// If error occurred, defer will close files; if success, resty closes them
	openedFiles = nil

	if err != nil {
		return nil, err
	}

	// Decode JSON response if output struct provided
	if out != nil {
		if unmarshalErr := json.Unmarshal(resp.Body(), out); unmarshalErr != nil {
			return resp, fmt.Errorf("failed to unmarshal response: %w", unmarshalErr)
		}
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithFields(logrus.Fields{
		"url":      req.URL,
		"method":   i.Method,
		"jobId":    i.JobID,
		"formData": structs.StructToJson(formData),
		"files":    len(files),
		"resBody":  resp.String(),
		"status":   resp.StatusCode(),
		"time":     time.Now(),
		"latency":  resp.Time(),
	}).Infoln("Outgoing Multipart Request")

	return resp, err
}
