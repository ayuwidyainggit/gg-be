package repository

import "master/model"

type EmployeeScopeRepository interface {
	FindEmployeeDropdownScope(empID int, custID string) (model.Employee, error)
}
