package service

import (
	"mobile/entity"
	"mobile/pkg/structs"
	"mobile/repository"
)

type OutletListService interface {
	List(filter entity.OutletListQueryFilter, custId string) ([]entity.OutletListItem, int64, int, error)
	Delete(custId string, outletId int64, deletedBy int64) error
	Update(request entity.UpdateOutletBody) error
}

type outletListServiceImpl struct {
	Repository repository.OutletListRepository
}

func NewOutletListService(repo repository.OutletListRepository) OutletListService {
	return &outletListServiceImpl{
		Repository: repo,
	}
}

// List - GET /v1/outlet-list
func (s *outletListServiceImpl) List(filter entity.OutletListQueryFilter, custId string) ([]entity.OutletListItem, int64, int, error) {
	outlets, total, lastPage, err := s.Repository.FindAll(filter, custId)
	if err != nil {
		return nil, 0, 0, err
	}

	var data []entity.OutletListItem
	for _, outlet := range outlets {
		var item entity.OutletListItem
		structs.Automapper(outlet, &item)
		data = append(data, item)
	}

	return data, total, lastPage, nil
}

// Delete - DELETE /v1/outlet-list/:outlet_id
func (s *outletListServiceImpl) Delete(custId string, outletId int64, deletedBy int64) error {
	return s.Repository.SoftDelete(custId, outletId, deletedBy)
}

// Update - PATCH /v1/m-outlets/:outlet_id
func (s *outletListServiceImpl) Update(request entity.UpdateOutletBody) error {
	// Build update map from non-nil fields
	updateData := map[string]interface{}{
		"updated_by": request.UpdatedBy,
	}

	if request.OutletName != nil {
		updateData["outlet_name"] = *request.OutletName
	}
	if request.Address != nil {
		updateData["address1"] = *request.Address
	}
	if request.PhoneNo != nil {
		updateData["phone_no"] = *request.PhoneNo
	}
	if request.BuildingOwn != nil {
		updateData["building_own"] = *request.BuildingOwn
	}
	if request.Latitude != nil {
		updateData["latitude"] = *request.Latitude
	}
	if request.Longitude != nil {
		updateData["longitude"] = *request.Longitude
	}
	if request.FileUrl != nil {
		updateData["file_url"] = *request.FileUrl
	}

	// Update outlet main data
	err := s.Repository.Update(request.CustId, request.OutletId, updateData)
	if err != nil {
		return err
	}

	// Update contact details if provided
	if request.Details != nil && len(request.Details.Contact) > 0 {
		for _, contact := range request.Details.Contact {
			err := s.Repository.UpdateContact(request.CustId, request.OutletId, contact)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
