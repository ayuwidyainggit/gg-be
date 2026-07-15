package responsebuild

import (
	"inventory/pkg/texttranslator"
	"log"
)

type DataRespReq struct {
	Translator *texttranslator.Translator
	Message    string
	Data       interface{}
	Errors     interface{}
	Paging     interface{}
	Filter     interface{}
	RequestId  string
	Lang       string
}

type ApiPayload struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors,omitempty"`
	Paging    interface{} `json:"paging,omitempty"`
	Filter    interface{} `json:"filter,omitempty"`
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
		log.Println("err:", err.Error())
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
		Filter:    resp.Filter,
		RequestId: resp.RequestId,
	}
}

func (resp *DataRespReq) SetFilter(filter interface{}) {
	if filter != nil {
		resp.Filter = filter
	}
}
