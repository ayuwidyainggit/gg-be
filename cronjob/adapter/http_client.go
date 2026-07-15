package adapter

import (
	"cronjob/pkg/structs"
	"errors"
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

func (i *HttpClientInfo) Dispatch() (*resty.Response, error) {
	var (
		err  error
		resp *resty.Response
	)

	req := i.prepare()
	if i.Method == "" {
		return nil, errors.New("missing http method for dispatch request")
	}

	if i.Payload == nil {
		resp, err = req.Execute(i.Method, i.Url)
	} else {
		resp, err = req.SetBody(i.Payload).Execute(i.Method, i.Url)
	}

	if err != nil {
		return nil, err
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
