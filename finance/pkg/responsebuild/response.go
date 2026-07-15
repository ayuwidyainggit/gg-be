package responsebuild

import (
	"finance/pkg/texttranslator"
	"github.com/gofiber/fiber/v2/log"
)

type DataRespReq struct {
	Translator *texttranslator.Translator
	Message    string
	Data       interface{}
	Errors     interface{}
	Paging     interface{}
	RequestId  string
	Lang       string
}

type ApiPayload struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Paging    interface{} `json:"paging,omitempty"`
	RequestId string      `json:"request_id"`
}

func BuildResponse(requestID string, lang ...string) *DataRespReq {
	// load translator
	Lang := texttranslator.EN
	if len(lang) > 0 {
		if lang[0] != "" {
			if lang[0] != texttranslator.EN {
				Lang = texttranslator.ID
			}
		}

	}
	trans, err := texttranslator.New(1)
	if err != nil {
		log.Error("err:", err.Error())
		panic(err)
	}
	return &DataRespReq{
		Translator: trans,
		RequestId:  requestID,
		Lang:       Lang,
	}
}

func (resp *DataRespReq) Setmsg(msg string) {
	resp.Message = msg
}
func (resp *DataRespReq) Setdata(data interface{}) {
	if data != nil {
		resp.Data = data
	}

}
func (resp *DataRespReq) Seterrors(errors interface{}) {
	if errors != nil {
		resp.Errors = errors
	}

}

func (resp *DataRespReq) Setpaging(paging interface{}) {
	if paging != nil {
		resp.Paging = paging
	}
}

func (resp *DataRespReq) GetRespPayload() *ApiPayload {
	if resp.Lang == texttranslator.ID && resp.Message != "" {
		resp.Message = resp.Translator.Translate(resp.Lang, resp.Message)
	}
	return &ApiPayload{
		Message:   resp.Message,
		Data:      resp.Data,
		Errors:    resp.Errors,
		Paging:    resp.Paging,
		RequestId: resp.RequestId,
	}
}
