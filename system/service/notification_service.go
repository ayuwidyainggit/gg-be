package service

import (
	"encoding/json"
	"errors"
	"system/adapter"
	"system/entity"
	"system/pkg/config/env"

	"github.com/gofiber/fiber/v2/log"
)

type NotificationService interface {
	WhatsappCicd(request entity.NotifyCicdWaReq) (resp entity.NotifyWaRes, err error)
}

type NotificationServiceImpl struct {
	Config     env.ConfigEnv
	HttpClient adapter.HttpClientInfo
}

func NewNotificationService(
	config env.ConfigEnv,
	httpClient adapter.HttpClientInfo,
) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		Config:     config,
		HttpClient: httpClient,
	}
}

func (service *NotificationServiceImpl) WhatsappCicd(req entity.NotifyCicdWaReq) (resp entity.NotifyWaRes, err error) {

	createBodyMsg :=
		`*Deployment to ` + req.Env + ` Environment has been COMPLETED with details :*
	
	Service Name : ` + req.ProjectName + `
	Branch : ` + req.Branch + `
	Commit Author : ` + req.CommitAuthor + `
	Commit ID : ` + req.CommitID + `
	Commit Message : ` + req.CommitMessage + `
	
	Developers please check the latest version of this service. Thank You.`

	request := adapter.HttpClientInfo{
		Url:    "https://gate.whapi.cloud/messages/text",
		Method: "POST",
		Auth:   "Bearer 3UehiKpfEMnxfr2Uq0Si3d7vqTH4md8C",
		Payload: map[string]interface{}{
			"to":   req.To,
			"body": createBodyMsg,
		},
	}

	response, err := request.Dispatch()
	if err != nil {
		log.Error("WhatsappCicd, Dispatch, error: ", err)
		return resp, err
	}

	if response.StatusCode() < 200 {
		errResp := entity.NotifyWaErrRes{}
		err = json.Unmarshal(response.Body(), &errResp)
		if err != nil {
			log.Error("unmarshal WhatsappCicd error response:", err)
			return resp, err
		}
		return resp, errors.New(errResp.ErrRes.Message)

	}

	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		log.Error("Error unmarshal WhatsappCicd success response:", err)
		return resp, err
	}

	return
}
