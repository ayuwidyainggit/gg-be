package service

import (
	"errors"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/rabbitmq"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type ConsProductService interface {
	Detail(int, string) (entity.ConsProductResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.ConsProductListResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.ConsProductLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateConsProductBody) (entity.ConsProductResponse, error)
	Update(int, entity.UpdateConsProductRequest) error
	Delete(string, int, int64) error
}

func NewConsProductService(consProductRepository repository.ConsProductRepository, mProductRepository repository.ProductRepository) *consProductServiceImpl {
	return &consProductServiceImpl{
		ConsProductRepository: consProductRepository,
		ProductRepository:     mProductRepository,
	}
}

type consProductServiceImpl struct {
	ConsProductRepository repository.ConsProductRepository
	ProductRepository     repository.ProductRepository
}

func (service *consProductServiceImpl) Detail(cProId int, custId string) (response entity.ConsProductResponse, err error) {
	consProduct, err := service.ConsProductRepository.FindOneByCProIdAndCustId(cProId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(consProduct, &response)
	if err != nil {
		return response, err
	}

	rmqConfig := rabbitmq.RmqConfig{
		ExchangeName:   "events",
		RoutingKey:     custId + ".products.events.detail",
		QueueName:      custId + ".products.events.detail",
		DelayQueueName: custId + ".products.events.detail.delay",
		MessageTTL:     "20000",
		Message:        structs.StructToJson(response),
	}
	err = rabbitmq.PublishMessage(&rmqConfig)
	if err != nil {
		log.Println("err >>>", err.Error())
	}

	return response, err
}

func (service *consProductServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ConsProductListResponse, total int, lastPage int, err error) {
	consProducts, total, lastPage, err := service.ConsProductRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range consProducts {
		var vResp entity.ConsProductListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	// err = rabbitmq.PublishMessage("events", custId+".products.events.list", custId+".products.events.*list", structs.StructToJson(data))
	// if err != nil {
	// 	log.Println("err >>>", err.Error())
	// }

	return data, total, lastPage, err
}

func (service *consProductServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ConsProductLookupResponse, total int, lastPage int, err error) {
	consProducts, total, lastPage, err := service.ConsProductRepository.FindAllLookupByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range consProducts {
		var vResp entity.ConsProductLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	// go rabbitmq.Subscribe(custId + ".products.events.list")   // Receives products.events  messages
	// go rabbitmq.Subscribe(custId + ".products.events.detail") // Receives products.events  messages

	return data, total, lastPage, err
}

func (service *consProductServiceImpl) Store(request entity.CreateConsProductBody) (response entity.ConsProductResponse, err error) {

	// c_pro_code & cust id validation, if err == nil, this means that code & cust id already exists
	consProduct, err := service.ConsProductRepository.FindOneByCProCodeAndCustId(request.CProCode, request.CustId)
	if err == nil {
		return response, errors.New("c_pro_code: " + consProduct.CProCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	consPro := model.ConsProduct{
		CustId:    request.CustId,
		CProCode:  request.CProCode,
		CProName:  request.CProName,
		IsActive:  request.IsActive,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
		UpdatedBy: &request.CreatedBy,
	}

	cProId, err := service.ConsProductRepository.Store(consPro)
	if err != nil {
		return response, err
	}

	response.CProId = cProId

	return response, err
}

func (service *consProductServiceImpl) Update(cProId int, request entity.UpdateConsProductRequest) (err error) {

	// c_pro_code & cust id validation, if err == nil and params cProId != consProduct.Id, this means that code & cust id already exists
	consProduct, err := service.ConsProductRepository.FindOneByCProCodeAndCustId(request.CProCode, request.CustId)
	if err == nil && consProduct.CProId != cProId {
		return errors.New("c_pro_code: " + consProduct.CProCode + " is already exists")
	}

	err = service.ConsProductRepository.Update(cProId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *consProductServiceImpl) Delete(custId string, cProId int, userId int64) (err error) {

	isExists, err := service.ProductRepository.IsExists(cProId, custId, "c_pro_id")
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("c_pro_id is still being used")
	}

	err = service.ConsProductRepository.Delete(custId, cProId, userId)
	if err != nil {
		return err
	}

	return err
}
