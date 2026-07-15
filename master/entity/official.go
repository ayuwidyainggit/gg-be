package entity

import (
	"time"
)

type OfficialQueryFilter struct {
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	OfficialType int    `query:"official_type"`
}

type OfficialResponse struct {
	OfficialId      int     `json:"official_id"`
	OfficialType    int     `json:"official_type"`
	OfficialName    *string `json:"official_name"`
	EmpId           int     `json:"emp_id"`
	EmpCode         string  `json:"emp_code"`
	EmpName         string  `json:"emp_name"`
	HierarchyCode   string  `json:"hierarchy_code"`
	SupervisorId2   *int    `json:"supervisor_id2"`
	SupervisorCode2 *string `json:"supervisor_code2"`
	SupervisorName2 *string `json:"supervisor_name2"`
	HierarchyCode2  string  `json:"hierarchy_code2"`
	SupervisorId1   *int    `json:"supervisor_id1"`
	SupervisorCode1 *string `json:"supervisor_code1"`
	SupervisorName1 *string `json:"supervisor_name1"`
	HierarchyCode1  string  `json:"hierarchy_code1"`
}

type OfficialListResponse struct {
	OfficialId      int        `json:"official_id"`
	OfficialType    int        `json:"official_type"`
	OfficialName    *string    `json:"official_name"`
	EmpId           int        `json:"emp_id"`
	EmpCode         string     `json:"emp_code"`
	EmpName         string     `json:"emp_name"`
	HierarchyCode   string     `json:"hierarchy_code"`
	SupervisorId2   *int       `json:"supervisor_id2"`
	SupervisorCode2 *string    `json:"supervisor_code2"`
	SupervisorName2 *string    `json:"supervisor_name2"`
	HierarchyCode2  string     `json:"hierarchy_code2"`
	SupervisorId1   *int       `json:"supervisor_id1"`
	SupervisorCode1 *string    `json:"supervisor_code1"`
	SupervisorName1 *string    `json:"supervisor_name1"`
	HierarchyCode1  string     `json:"hierarchy_code1"`
	CreatedAt       *time.Time `json:"created_at"`
}

type OfficialLookupResponse struct {
	OfficialId      int     `json:"official_id"`
	OfficialType    int     `json:"official_type"`
	OfficialName    *string `json:"official_name"`
	EmpId           int     `json:"emp_id"`
	EmpCode         string  `json:"emp_code"`
	EmpName         string  `json:"emp_name"`
	HierarchyCode   string  `json:"hierarchy_code"`
	SupervisorId2   *int    `json:"supervisor_id2"`
	SupervisorCode2 *string `json:"supervisor_code2"`
	SupervisorName2 *string `json:"supervisor_name2"`
	HierarchyCode2  string  `json:"hierarchy_code2"`
	SupervisorId1   *int    `json:"supervisor_id1"`
	SupervisorCode1 *string `json:"supervisor_code1"`
	SupervisorName1 *string `json:"supervisor_name1"`
	HierarchyCode1  string  `json:"hierarchy_code1"`
}

type CreateOfficialBody struct {
	CustId    string               `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64                `json:"created_by" validate:"required"`
	Officials []CreateOfficialData `json:"officials" validate:"required"`
}

type CreateOfficialData struct {
	OfficialType int `json:"official_type" validate:"required,max=3"`
	EmpId        int `json:"emp_id" validate:"required"`
	SupervisorId int `json:"supervisor_id" validate:""`
}

type CreateOfficialBodyHierarchy struct {
	CustId    string                        `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64                         `json:"created_by" validate:"required"`
	Officials []CreateOfficialDataHierarchy `json:"officials" validate:"required"`
}

type CreateOfficialDataHierarchy struct {
	OfficialType int                           `json:"official_type" validate:"required,max=3"`
	EmpId        int                           `json:"emp_id" validate:"required,unique,dive"`
	SupervisorId int                           `json:"supervisor_id" validate:""`
	Children     []CreateOfficialDataHierarchy `json:"children" validate:""`
}

type DetailOfficialParams struct {
	OfficialId int `params:"official_id" validate:"required"`
}

type UpdateOfficialParams struct {
	OfficialId int `params:"official_id" validate:"required"`
}

type DeleteOfficialParams struct {
	OfficialId int `params:"official_id" validate:"required"`
}

type UpdateOfficialRequest struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy    int64  `json:"updated_by" validate:"required"`
	OfficialType int    `json:"official_type,omitempty" validate:"max=3,omitempty"`
	OfficialName string `json:"official_name,omitempty" validate:"max=150,omitempty"`
	SupervisorId int    `json:"supervisor_id,omitempty" validate:"omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

type OfficialHierarchyResp struct {
	OfficialId     int                     `json:"official_id"`
	OfficialType   int                     `json:"official_type"`
	OfficialName   *string                 `json:"official_name"`
	EmpId          int                     `json:"emp_id"`
	EmpName        string                  `json:"emp_name"`
	HierarchyCode  string                  `json:"hierarchy_code"`
	SupervisorId   int                     `json:"supervisor_id"`
	SupervisorName string                  `json:"supervisor_name"`
	Children       []OfficialHierarchyResp `json:"children,omitempty"`
}

type OfficialHierarchyMap struct {
	Db   map[int][]OfficialHierarchyResp
	Resp []OfficialHierarchyResp
}

func NewOfficialHierarchyMap() *OfficialHierarchyMap {
	return &OfficialHierarchyMap{
		Db: make(map[int][]OfficialHierarchyResp),
	}
}
func (w *OfficialHierarchyMap) SetChildrenRecursively(res *OfficialHierarchyResp) {
	// append to make a copy of the slice (otherwise we will be changing items in the 'database')
	res.Children = append([]OfficialHierarchyResp{}, w.Db[res.EmpId]...) // Get the children from simulated database

	for i := range res.Children {
		w.SetChildrenRecursively(&res.Children[i])
	}
}
func (w *OfficialHierarchyMap) Append(resp OfficialHierarchyResp) {
	w.Resp = append(w.Resp, resp)
}
