package responsebuild

type DataRespReq struct {
	Message   string
	Data      interface{}
	Errors    interface{}
	Paging    interface{}
	RequestId string
	Lang      string
}

type ApiPayload struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Paging    interface{} `json:"paging,omitempty"`
	RequestId string      `json:"request_id"`
}

func BuildResponse(requestID string, lang ...string) *DataRespReq {
	return &DataRespReq{
		RequestId: requestID,
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
	return &ApiPayload{
		Message:   resp.Message,
		Data:      resp.Data,
		Errors:    resp.Errors,
		Paging:    resp.Paging,
		RequestId: resp.RequestId,
	}
}
