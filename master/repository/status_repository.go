package repository

import (
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type StatusRepository interface {
	FindOneByStatusIdAndStatusValue(statusId string, statusVal int, langId string) (model.Status, error)
	FindAllByStatusIdAndLangId(filter entity.StatusQueryFilter) (statuses []model.Status, total int, lastPage int, err error)
	FindAllByStatusIdAndLangIdLookup(filter entity.StatusQueryFilter) (statuses []model.Status, total int, lastPage int, err error)
	// Store(Status model.Status) (string, error)
	// Update(StatusId int, request entity.UpdateStatusRequest) error
	// Delete(custId string, StatusId int, deletedBy int64) error
}

func NewStatusRepository(db *sqlx.DB) StatusRepository {
	return &statusRepositoryImpl{db}
}

type statusRepositoryImpl struct {
	*sqlx.DB
}

func (repository *statusRepositoryImpl) FindOneByStatusIdAndStatusValue(statusId string, statusVal int, langId string) (model.Status, error) {
	status := model.Status{}
	query := `SELECT 
				status_id, status_name, status_value, lang_id
			FROM mst.m_status
			WHERE status_id = $1 AND status_value = $2 AND lang_id = $3`
	err := repository.Get(&status, query, statusId, statusVal, langId)
	if err != nil {
		log.Println("StatusRepository, FindOneByStatusTypeAndCustId, err:", err.Error())
		return status, err
	}

	return status, nil
}

func (repository *statusRepositoryImpl) FindAllByStatusIdAndLangId(filter entity.StatusQueryFilter) ([]model.Status, int, int, error) {

	Statuss := []model.Status{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` status_id, status_name, status_value, lang_id`
	qWhere := ` WHERE 1=1 `

	if filter.StatusId != "" {
		qWhere += ` AND status_id = '` + filter.StatusId + `' `
	}

	if filter.LangId != "" {
		qWhere += ` AND lang_id = '` + filter.LangId + `' `
	}

	qFrom := ` FROM mst.m_status `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("StatusRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("StatusRepository, count total, err:", err.Error())
		return Statuss, 0, 0, err
	}

	sortBy := `` // default sort by
	if filter.Sort != "" {
		mSortBy := strings.Split(filter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `status_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if filter.Limit == 0 {
		filter.Limit = 10
	}

	page := filter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * filter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(filter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(filter.Limit), strconv.Itoa(offset))

	// log.Println("StatusRepository, querySelect:", querySelect)
	err = repository.Select(&Statuss, querySelect)
	if err != nil {
		log.Println("StatusRepository, FindAllByCustId, err:", err.Error())
		return Statuss, total, lastPage, err
	}

	return Statuss, total, lastPage, nil
}

func (repository *statusRepositoryImpl) FindAllByStatusIdAndLangIdLookup(filter entity.StatusQueryFilter) ([]model.Status, int, int, error) {

	Statuss := []model.Status{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` status_id, status_name, status_value, lang_id `
	qWhere := ` WHERE 1=1 `

	if filter.StatusId != "" {
		qWhere += ` AND status_id = '` + filter.StatusId + `' `
	}

	if filter.LangId != "" {
		qWhere += ` AND lang_id = '` + filter.LangId + `' `
	}

	qFrom := ` FROM mst.m_status `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("StatusRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("StatusRepository, count total, err:", err.Error())
		return Statuss, 0, 0, err
	}

	sortBy := `` // default sort by
	if filter.Sort != "" {
		mSortBy := strings.Split(filter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `status_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	lastPage := 1

	// log.Println("StatusRepository, querySelect:", querySelect)
	err = repository.Select(&Statuss, querySelect)
	if err != nil {
		log.Println("StatusRepository, FindAllByCustIdLookup, err:", err.Error())
		return Statuss, total, lastPage, err
	}

	return Statuss, total, lastPage, nil
}

// func (repository *statusRepositoryImpl) Store(Status model.Status) (string, error) {
// 	query :=
// 		`INSERT INTO mst.m_status(status_id, status_name, status_value, lang_id)
// 		VALUES ($1, $2, $3, $4) RETURNING status_id;`
// 	lastInsertId := Status.StatusId
// 	err := repository.QueryRow(query, Status.StatusId, Status.StatusName, Status.LangId).Scan(&lastInsertId)
// 	if err != nil {
// 		log.Println("StatusRepository, Store, err:", err.Error())
// 		return Status.StatusId, err
// 	}
// 	return Status.StatusId, nil
// }

// func (repository *statusRepositoryImpl) Update(statusId string, request entity.UpdateStatusRequest) error {
// 	var (
// 		r            model.StatusUpdate
// 		sqlSetFields string
// 		nRows        int64
// 	)

// 	reqByte, _ := json.Marshal(request)
// 	_ = json.Unmarshal(reqByte, &r)
// 	sqlPatch := sql_helper.SQLPatches(r)

// 	// data, _ := json.Marshal(sqlPatch)
// 	// fmt.Printf("StatusRepository, Update, Fields & Args: %s\n", data)

// 	for i, _ := range sqlPatch.Fields {
// 		sqlSetFields += sqlPatch.Fields[i] + ", "
// 		i++
// 	}
// 	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

// 	query := `UPDATE mst.m_status
// 			  SET ` + sqlSetFields + `,
// 			  	updated_at = CURRENT_TIMESTAMP
// 			  WHERE is_del = false
// 			  AND cust_id = :cust_id
// 			  AND status_id = :status_id_old;`

// 	// log.Println("StatusRepository, Update, query:", query)

// 	sqlPatch.Args["status_id_old"] = StatusId
// 	sqlPatch.Args["cust_id"] = requeCustId

// 	result, err := repository.NamedExec(query, sqlPatch.Args)
// 	if err != nil {
// 		log.Println("StatusRepository, Update, err:", err.Error())
// 		return err
// 	}

// 	if nRows, err = result.RowsAffected(); err != nil {
// 		return errors.New("no rows affected")
// 	}
// 	if nRows == 0 {
// 		return errors.New("no rows affected")
// 	}

// 	return nil
// }

// func (repository *statusRepositoryImpl) Delete(custId string, StatusId int, deletedBy int64) error {
// 	var nRows int64
// 	query := `UPDATE mst.m_status
// 			SET is_del = true,
// 				deleted_at = CURRENT_TIMESTAMP,
// 				deleted_by = :deleted_by
// 			WHERE is_del = false
// 			AND cust_id = :cust_id
// 			AND status_id = :status_id;`

// 	wMap := map[string]interface{}{
// 		"cust_id":    custId,
// 		"status_id":  StatusId,
// 		"deleted_by": deletedBy,
// 	}

// 	result, err := repository.NamedExec(query, wMap)
// 	if err != nil {
// 		log.Println("StatusRepository, Delete, err:", err.Error())
// 		return err
// 	}

// 	if nRows, err = result.RowsAffected(); err != nil {
// 		return errors.New("no rows affected")
// 	}
// 	if nRows == 0 {
// 		return errors.New("no rows affected")
// 	}

// 	return nil
// }
