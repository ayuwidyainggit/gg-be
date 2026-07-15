package model

import (
	"time"
)

type Official struct {
	CustId       string     `db:"cust_id" json:"cust_id"`
	OfficialId   int        `db:"official_id" json:"official_id"`
	OfficialType int        `db:"official_type" json:"official_type"`
	OfficialName *string    `db:"official_name" json:"official_name"`
	EmpId        *int       `db:"emp_id" json:"emp_id"`
	SupervisorId *int       `db:"supervisor_id" json:"supervisor_id"`
	CreatedBy    *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt    *time.Time `db:"created_at,omitempty" json:"created_at"`
}

type AllOfficialHierarchy struct {
	CustId         string     `db:"cust_id" json:"cust_id"`
	OfficialId     int        `db:"official_id" json:"official_id"`
	OfficialType   int        `db:"official_type" json:"official_type"`
	OfficialName   *string    `db:"official_name" json:"official_name"`
	HierarchyCode  *string    `db:"hierarchy_code" json:"hierarchy_code"`
	EmpId          *int       `db:"emp_id" json:"emp_id"`
	EmpName        *string    `db:"emp_name" json:"emp_name"`
	SupervisorId   *int       `db:"supervisor_id" json:"supervisor_id"`
	SupervisorName *string    `db:"supervisor_name" json:"supervisor_name"`
	CreatedBy      *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt      *time.Time `db:"created_at,omitempty" json:"created_at"`
}

type OfficialList struct {
	CustId          string     `db:"cust_id" json:"cust_id"`
	OfficialId      int        `db:"official_id" json:"official_id"`
	OfficialType    int        `db:"official_type" json:"official_type"`
	OfficialName    *string    `db:"official_name" json:"official_name"`
	EmpId           *int       `db:"emp_id" json:"emp_id"`
	EmpCode         *string    `db:"emp_code" json:"emp_code"`
	EmpName         *string    `db:"emp_name" json:"emp_name"`
	HierarchyCode   *string    `db:"hierarchy_code" json:"hierarchy_code"`
	SupervisorId2   *int       `db:"supervisor_id2" json:"supervisor_id2"`
	SupervisorCode2 *string    `db:"supervisor_code2" json:"supervisor_code2"`
	SupervisorName2 *string    `db:"supervisor_name2" json:"supervisor_name2"`
	HierarchyCode2  *string    `db:"hierarchy_code2" json:"hierarchy_code2"`
	SupervisorId1   *int       `db:"supervisor_id1" json:"supervisor_id1"`
	SupervisorName1 *string    `db:"supervisor_name1" json:"supervisor_name1"`
	SupervisorCode1 *string    `db:"supervisor_code1" json:"supervisor_code1"`
	HierarchyCode1  *string    `db:"hierarchy_code1" json:"hierarchy_code1"`
	CreatedBy       *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt       *time.Time `db:"created_at,omitempty" json:"created_at"`
}

type OfficialUpdate struct {
	OfficialType *int       `json:"official_type,omitempty" sql:"official_type"`
	OfficialName *string    `db:"official_name" json:"official_name"`
	EmpId        *int       `json:"emp_id,omitempty" sql:"emp_id"`
	SupervisorId *int       `json:"supervisor_id,omitempty" sql:"supervisor_id"`
	IsActive     *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt    *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy    *int64     `json:"updated_by" sql:"updated_by"`
}

type OfficialEmployee struct {
	CustId   string `db:"cust_id" json:"cust_id"`
	EmpId    int    `db:"emp_id" json:"emp_id"`
	EmpCode  string `db:"emp_code" json:"emp_code"`
	EmpName  string `db:"emp_name" json:"emp_name"`
	IsActive bool   `db:"is_active" json:"is_active"`
}
