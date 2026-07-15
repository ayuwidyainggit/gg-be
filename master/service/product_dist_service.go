package service

import (
	"master/entity"
	"master/pkg/structs"
	"master/repository"
)

type ProductDistService interface {
	Detail(int64, string, string) (entity.ProductDistResponse, error)
	List(entity.ProductDistQueryFilter, string, string) (data []entity.ProductDistResponse, total int, lastPage int, err error)
	LookupList(entity.ProductDistQueryFilter, string, string) (data []entity.ProductDistLookupResponse, total int, lastPage int, err error)
	SearchList(entity.ProductDistQueryFilter, string, string) (data []entity.ProductDistSearchResponse, total int, lastPage int, err error)
	Store(entity.CreateProductBody) (entity.ProductDistResponse, error)
	Update(int64, entity.UpdateProductDistRequest) error
	Delete(string, int64, int64) error
}

func NewProductDistService(ProductDistRepository repository.ProductDistRepository) *ProductDistServiceImpl {
	return &ProductDistServiceImpl{
		ProductDistRepository: ProductDistRepository,
	}
}

type ProductDistServiceImpl struct {
	ProductDistRepository repository.ProductDistRepository
}

func (service *ProductDistServiceImpl) Detail(productId int64, custId, parentCustId string) (response entity.ProductDistResponse, err error) {
	product, err := service.ProductDistRepository.FindOneByProductIdAndCustId(productId, custId, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(product, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *ProductDistServiceImpl) List(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) (data []entity.ProductDistResponse, total int, lastPage int, err error) {
	products, total, lastPage, err := service.ProductDistRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductDistResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ProductDistServiceImpl) LookupList(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) (data []entity.ProductDistLookupResponse, total int, lastPage int, err error) {
	products, total, lastPage, err := service.ProductDistRepository.FindAllByCustIdLookup(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductDistLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ProductDistServiceImpl) SearchList(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) (data []entity.ProductDistSearchResponse, total int, lastPage int, err error) {
	products, total, lastPage, err := service.ProductDistRepository.FindAllByCustIdSearch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductDistSearchResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ProductDistServiceImpl) Store(request entity.CreateProductBody) (response entity.ProductDistResponse, err error) {

	// product_code & cust id validation, if err == nil, this means that code & cust id already exists
	// product, err := service.ProductDistRepository.FindOneByProductCodeAndCustId(request.ProductCode, request.CustId)
	// if err == nil {
	// 	return response, errors.New("product_code: " + product.ProductCode + " is already exists")
	// }

	// timeNow := time.Now().In(time.UTC)
	// productData := model.Product{}
	// err = structs.Automapper(request, &productData)
	// if err != nil {
	// 	return response, err
	// }

	// productData.CustId = request.CustId
	// productData.ProductCode = request.ProductCode
	// productData.ProductName = request.ProductName
	// productData.IsActive = request.IsActive
	// productData.CreatedAt = &timeNow
	// productData.CreatedBy = &request.CreatedBy
	// productData.UpdatedAt = &timeNow
	// productData.UpdatedBy = &request.CreatedBy

	// productId, err := service.ProductDistRepository.Store(productData)
	// if err != nil {
	// 	return response, err
	// }

	// response.ProductId = productId

	return response, err
}

func (service *ProductDistServiceImpl) Update(productId int64, request entity.UpdateProductDistRequest) (err error) {

	err = service.ProductDistRepository.Update(productId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *ProductDistServiceImpl) Delete(custId string, productId int64, userId int64) (err error) {

	err = service.ProductDistRepository.Delete(custId, productId, userId)
	if err != nil {
		return err
	}

	return err
}
