package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type SAPReplStatusRequest struct {
	CustID          string                       `json:"cust_id"`
	DistributorCode int                          `json:"distributor_code"`
	Replenishments  []SAPReplStatusReplenishment `json:"replenishments" validate:"required,min=1"`
}

type SAPReplStatusReplenishment struct {
	ReplenishmentNo string                    `json:"replenishment_no" validate:"required"`
	Status          int                       `json:"status" validate:"required"`
	Details         []SAPReplStatusDetailLine `json:"details" validate:"required,min=1,dive"`
}

type SAPReplStatusDetailLine struct {
	ProCode     ProCodeJSON `json:"pro_code"`
	PurchPrice3 *float64    `json:"purch_price3,omitempty"`
	Qty3        *float64    `json:"qty3,omitempty"`
	UnitID3     string      `json:"unit_id3" validate:"required"`
}

type ProCodeJSON string

func (p ProCodeJSON) String() string { return strings.TrimSpace(string(p)) }

func (p *ProCodeJSON) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*p = ""
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*p = ProCodeJSON(strings.TrimSpace(s))
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		*p = ProCodeJSON(strings.TrimSpace(n.String()))
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("pro_code must be string or number")
	}
	*p = ProCodeJSON(strings.TrimSpace(s))
	return nil
}


type SAPReplStatusResponse struct {
	RequestID string                     `json:"request_id"`
	Status    string                     `json:"status"`
	Message   string                     `json:"message"`
	Data      interface{}                `json:"data,omitempty"`
	Errors    []SAPReplStatusReplErrWrap `json:"errors,omitempty"`
}

type SAPReplStatusOKItem struct {
	ReplenishmentNo string `json:"replenishment_no"`
	Status          string `json:"status"`
}

type SAPReplStatusPartialData struct {
	Success []SAPReplStatusOKItem     `json:"success"`
	Failed  []SAPReplStatusFailedItem `json:"failed"`
}

type SAPReplStatusFailedItem struct {
	ReplenishmentNo string              `json:"replenishment_no"`
	Status          string              `json:"status"`
	Errors          []SAPReplFieldError `json:"errors"`
}

type SAPReplFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type SAPReplStatusReplErrWrap struct {
	ReplenishmentNo string              `json:"replenishment_no"`
	Errors          []SAPReplFieldError `json:"errors"`
}
