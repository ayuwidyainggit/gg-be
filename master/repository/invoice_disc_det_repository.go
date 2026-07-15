package repository

import (
	"errors"
	"log"
	"master/model"
	"master/pkg/structs"

	"github.com/jmoiron/sqlx"
)

type InvoiceDiscDetRepository interface {
	FindOneByInvDiscIdAndCustId(invDiscId int, custId string) ([]model.InvoiceDiscDet, error)
	BulkInsert(invDiscDetails []model.InvoiceDiscDet) error
	DeleteByInvDiscId(invDiscId int, custId string) error
}

func NewInvoiceDiscDetRepository(db *sqlx.DB) InvoiceDiscDetRepository {
	return &InvoiceDiscDetRepositoryImpl{db}
}

type InvoiceDiscDetRepositoryImpl struct {
	*sqlx.DB
}

func (repository *InvoiceDiscDetRepositoryImpl) FindOneByInvDiscIdAndCustId(invDiscId int, custId string) ([]model.InvoiceDiscDet, error) {
	invoiceDiscDets := []model.InvoiceDiscDet{}
	query := `SELECT 
				cust_id, inv_disc_id, row_no,
				min_value, max_value, disc_perc
			  FROM mst.m_invoice_disc_det 
			  WHERE inv_disc_id = $1 
			  AND cust_id = $2 
			  ORDER BY row_no ASC`
	err := repository.Select(&invoiceDiscDets, query, invDiscId, custId)
	if err != nil {
		log.Println("InvoiceDiscDetRepository, FindOneByInvDiscIdAndCustId, err:", err.Error())
		return invoiceDiscDets, err
	}

	return invoiceDiscDets, nil
}

func (repository *InvoiceDiscDetRepositoryImpl) BulkInsert(invoiceDiscDetails []model.InvoiceDiscDet) error {
	invDiscDetails := make([]map[string]interface{}, 0)
	for _, row := range invoiceDiscDetails {
		mapDet := structs.Map(row)
		invDiscDetails = append(invDiscDetails, mapDet)
	}

	query :=
		`INSERT INTO mst.m_invoice_disc_det(
			cust_id, inv_disc_id, row_no, 
			min_value, max_value, disc_perc
		)
		VALUES ( 
			:cust_id, :inv_disc_id, :row_no, 
			:min_value, :max_value, :disc_perc
		)`

	_, err := repository.NamedExec(query, invDiscDetails)
	if err != nil {
		return err
	}

	return nil
}

func (repository *InvoiceDiscDetRepositoryImpl) DeleteByInvDiscId(invDiscId int, custId string) error {
	var nRows int64
	query := `DELETE FROM mst.m_invoice_disc_det
			  WHERE cust_id = :cust_id
			  AND inv_disc_id = :inv_disc_id`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"inv_disc_id": invDiscId,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("InvoiceDiscDetRepository, Delete, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
