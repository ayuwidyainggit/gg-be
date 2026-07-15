package repository

import (
	"fmt"
	"master/entity"
	"master/model"
	"math"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type SurveyRepository interface {
	FindAllByCustId(filter entity.SurveyQueryFilter, custId string) ([]model.Survey, int, int, error)
	FindOneById(surveyId int, custId string) (model.Survey, error)
	FindCustIdsByDistributorIds(parentCustId string, distributorIds []int) ([]string, error)
	FindSurveyAreasByDistributorIds(distributorIds []int) ([]model.SurveyArea, error)
	Store(tx *sqlx.Tx, survey model.Survey) (int, error)
	Update(tx *sqlx.Tx, surveyId int, custId string, survey model.Survey) error
	ExistsActiveTitleOverlap(custId, surveyTitle string, dateStart, dateEnd time.Time, excludeSurveyId *int) (bool, error)
	Deactivate(tx *sqlx.Tx, surveyId int, custId string, isActive bool, updatedBy int64) error
	// Survey Area
	StoreAreas(tx *sqlx.Tx, surveyAreas []model.SurveyArea) error
	DeleteAreasBySurveyId(tx *sqlx.Tx, surveyId int) error
	FindAreasBySurveyId(surveyId int) ([]model.SurveyArea, error)
	// Survey Salesman
	StoreSalesmen(tx *sqlx.Tx, surveySalesmen []model.SurveySalesman) error
	DeleteSalesmenBySurveyId(tx *sqlx.Tx, surveyId int) error
	FindSalesmenBySurveyId(surveyId int) ([]model.SurveySalesman, error)
	// Survey Outlet
	StoreOutlets(tx *sqlx.Tx, surveyId int, outletIds []int) error
	DeleteOutletsBySurveyId(tx *sqlx.Tx, surveyId int) error
	FindOutletsBySurveyId(surveyId int) ([]model.SurveyOutlet, error)
	// Survey Detail (Template mapping)
	StoreDetails(tx *sqlx.Tx, surveyId int, templateIds []int) error
	DeleteDetailsBySurveyId(tx *sqlx.Tx, surveyId int) error
	FindDetailsBySurveyId(surveyId int) ([]model.SurveyDetail, error)
	// Survey Distributor (Sprint 13 SX-2445/SX-2448/SX-2452)
	StoreSurveyDistributors(tx *sqlx.Tx, distributorIds []int, surveyId int, custId string, createdBy int64) error
	DeleteSurveyDistributorsBySurveyId(tx *sqlx.Tx, surveyId int) error
	FindSurveyDistributorsBySurveyId(surveyId int) ([]model.SurveyDistributor, error)
	BeginTx() (*sqlx.Tx, error)
}

func NewSurveyRepository(db *sqlx.DB) SurveyRepository {
	return &surveyRepositoryImpl{DB: db}
}

type surveyRepositoryImpl struct {
	*sqlx.DB
}

func (r *surveyRepositoryImpl) BeginTx() (*sqlx.Tx, error) {
	return r.DB.Beginx()
}

func (r *surveyRepositoryImpl) FindAllByCustId(filter entity.SurveyQueryFilter, custId string) ([]model.Survey, int, int, error) {
	var surveys []model.Survey
	var total int

	baseQuery := `FROM mst.m_survey WHERE cust_id = $1 AND is_del = false`
	args := []interface{}{custId}
	argIdx := 2

	// Filter by status
	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}

	// Filter by answer_frequency
	if filter.AnswerFrequency != nil && *filter.AnswerFrequency != "" {
		baseQuery += fmt.Sprintf(" AND answer_frequency = $%d", argIdx)
		args = append(args, *filter.AnswerFrequency)
		argIdx++
	}

	// Filter by response_type (supports both response_frequency and response_type query params)
	responseTypeFilter := filter.ResponseFrequency
	if responseTypeFilter == nil || *responseTypeFilter == "" {
		responseTypeFilter = filter.ResponseType
	}
	if responseTypeFilter != nil && *responseTypeFilter != "" {
		baseQuery += fmt.Sprintf(" AND response_type = $%d", argIdx)
		args = append(args, *responseTypeFilter)
		argIdx++
	}

	// Search filter
	if filter.Query != "" {
		baseQuery += fmt.Sprintf(" AND survey_title ILIKE $%d", argIdx)
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	// Count
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.DB.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, 0, err
	}

	// Pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit
	lastPage := int(math.Ceil(float64(total) / float64(filter.Limit)))

	// Sort
	sortField := "created_at"
	sortOrder := "DESC"
	if filter.Sort != "" {
		parts := strings.Split(filter.Sort, ":")
		if len(parts) == 2 {
			sortField = parts[0]
			sortOrder = strings.ToUpper(parts[1])
		}
	}

	dataQuery := fmt.Sprintf(`SELECT survey_id, cust_id, survey_title, answer_frequency, response_type, 
		target_type, emp_id, efective_date_start, efective_date_end, status, created_at 
		%s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, sortField, sortOrder, argIdx, argIdx+1)
	args = append(args, filter.Limit, offset)

	err = r.DB.Select(&surveys, dataQuery, args...)
	if err != nil {
		return nil, 0, 0, err
	}

	return surveys, total, lastPage, nil
}

func (r *surveyRepositoryImpl) FindOneById(surveyId int, custId string) (model.Survey, error) {
	var survey model.Survey
	query := `SELECT s.survey_id, s.cust_id, s.survey_title, s.answer_frequency, s.response_type,
		s.target_type, s.level_target, s.emp_id, s.efective_date_start, s.efective_date_end, s.status, s.created_at,
		sm.sales_name,
		md.distributor_id,
		md.distributor_code,
		md.distributor_name
		FROM mst.m_survey s
		LEFT JOIN mst.m_salesman sm ON s.emp_id = sm.emp_id AND sm.cust_id = s.cust_id AND sm.is_del = false
		LEFT JOIN (
			SELECT survey_id, MIN(distributor_id) AS distributor_id
			FROM mst.m_survey_area
			WHERE is_del = false
			GROUP BY survey_id
		) sa ON sa.survey_id = s.survey_id
		LEFT JOIN smc.m_customer mc ON mc.cust_id = s.cust_id
		LEFT JOIN mst.m_distributor md ON md.distributor_id = COALESCE(sa.distributor_id, mc.distributor_id) AND md.is_del = false
		WHERE s.survey_id = $1 AND s.cust_id = $2 AND s.is_del = false`
	err := r.DB.Get(&survey, query, surveyId, custId)
	return survey, err
}

func (r *surveyRepositoryImpl) FindCustIdsByDistributorIds(parentCustId string, distributorIds []int) ([]string, error) {
	if len(distributorIds) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`SELECT cust_id FROM smc.m_customer WHERE distributor_id IN (?) AND parent_cust_id = ?`, distributorIds, parentCustId)
	if err != nil {
		return nil, err
	}

	query = r.DB.Rebind(query)
	var custIds []string
	err = r.DB.Select(&custIds, query, args...)
	if err != nil {
		return nil, err
	}

	return custIds, nil
}

func (r *surveyRepositoryImpl) FindSurveyAreasByDistributorIds(distributorIds []int) ([]model.SurveyArea, error) {
	if len(distributorIds) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`SELECT distributor_id, area_id
		FROM mst.m_distributor
		WHERE distributor_id IN (?)
			AND area_id IS NOT NULL
			AND area_id > 0
			AND is_del = false
		ORDER BY distributor_id ASC`, distributorIds)
	if err != nil {
		return nil, err
	}

	query = r.DB.Rebind(query)
	var surveyAreas []model.SurveyArea
	err = r.DB.Select(&surveyAreas, query, args...)
	if err != nil {
		return nil, err
	}

	return surveyAreas, nil
}

func (r *surveyRepositoryImpl) Store(tx *sqlx.Tx, survey model.Survey) (int, error) {
	var id int
	query := `INSERT INTO mst.m_survey
		(cust_id, survey_title, answer_frequency, response_type, target_type, level_target, emp_id,
		efective_date_start, efective_date_end, status, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING survey_id`
	err := tx.Get(&id, query,
		survey.CustId, survey.SurveyTitle, survey.AnswerFrequency, survey.ResponseType,
		survey.TargetType, survey.LevelTarget, survey.EmpId, survey.EfectiveDateStart, survey.EfectiveDateEnd,
		survey.Status, survey.CreatedAt, survey.CreatedBy, survey.UpdatedAt, survey.UpdatedBy)
	return id, err
}

func (r *surveyRepositoryImpl) Update(tx *sqlx.Tx, surveyId int, custId string, survey model.Survey) error {
	query := `UPDATE mst.m_survey
		SET survey_title = $1, answer_frequency = $2, response_type = $3, target_type = $4,
		level_target = $5, emp_id = $6, efective_date_start = $7, efective_date_end = $8,
		updated_at = $9, updated_by = $10
		WHERE survey_id = $11 AND cust_id = $12 AND is_del = false`
	_, err := tx.Exec(query, survey.SurveyTitle, survey.AnswerFrequency, survey.ResponseType,
		survey.TargetType, survey.LevelTarget, survey.EmpId, survey.EfectiveDateStart, survey.EfectiveDateEnd,
		survey.UpdatedAt, survey.UpdatedBy, surveyId, custId)
	return err
}

func (r *surveyRepositoryImpl) ExistsActiveTitleOverlap(custId, surveyTitle string, dateStart, dateEnd time.Time, excludeSurveyId *int) (bool, error) {
	var total int

	query := `SELECT COUNT(*)
		FROM mst.m_survey
		WHERE cust_id = $1
			AND status = 1
			AND is_del = false
			AND LOWER(TRIM(survey_title)) = LOWER(TRIM($2))
			AND efective_date_start <= $3
			AND efective_date_end >= $4`

	args := []interface{}{custId, surveyTitle, dateEnd, dateStart}
	if excludeSurveyId != nil {
		query += " AND survey_id <> $5"
		args = append(args, *excludeSurveyId)
	}

	err := r.DB.Get(&total, query, args...)
	if err != nil {
		return false, err
	}

	return total > 0, nil
}

func (r *surveyRepositoryImpl) Deactivate(tx *sqlx.Tx, surveyId int, custId string, isActive bool, updatedBy int64) error {
	now := time.Now().UTC()
	status := 0
	if isActive {
		status = 1
	}
	query := `UPDATE mst.m_survey SET status = $1, updated_at = $2, updated_by = $3
		WHERE survey_id = $4 AND cust_id = $5 AND is_del = false`
	result, err := tx.Exec(query, status, now, updatedBy, surveyId, custId)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("survey not found")
	}
	return nil
}

// Survey Area methods
func (r *surveyRepositoryImpl) StoreAreas(tx *sqlx.Tx, surveyAreas []model.SurveyArea) error {
	for _, surveyArea := range surveyAreas {
		var query string
		var err error
		if surveyArea.TargetCustId != "" {
			query = `INSERT INTO mst.m_survey_area (survey_id, distributor_id, area_id, target_cust_id) VALUES ($1, $2, $3, $4)`
			_, err = tx.Exec(query, surveyArea.SurveyId, surveyArea.DistributorId, surveyArea.AreaId, surveyArea.TargetCustId)
		} else {
			query = `INSERT INTO mst.m_survey_area (survey_id, distributor_id, area_id) VALUES ($1, $2, $3)`
			_, err = tx.Exec(query, surveyArea.SurveyId, surveyArea.DistributorId, surveyArea.AreaId)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *surveyRepositoryImpl) DeleteAreasBySurveyId(tx *sqlx.Tx, surveyId int) error {
	query := `UPDATE mst.m_survey_area SET is_del = true WHERE survey_id = $1`
	_, err := tx.Exec(query, surveyId)
	return err
}

func (r *surveyRepositoryImpl) FindAreasBySurveyId(surveyId int) ([]model.SurveyArea, error) {
	var areas []model.SurveyArea
	query := `SELECT sa.survey_area_id, sa.survey_id, sa.distributor_id, sa.area_id, sa.target_cust_id,
		a.area_name, md.distributor_name, mc.cust_name
		FROM mst.m_survey_area sa
		LEFT JOIN mst.m_area a ON sa.area_id = a.area_id
		LEFT JOIN mst.m_distributor md ON md.distributor_id = sa.distributor_id AND md.is_del = false
		LEFT JOIN smc.m_customer mc ON mc.cust_id = NULLIF(sa.target_cust_id, '')
		WHERE sa.survey_id = $1 AND sa.is_del = false
		ORDER BY sa.distributor_id ASC, sa.area_id ASC, sa.survey_area_id ASC`
	err := r.DB.Select(&areas, query, surveyId)
	return areas, err
}

// Survey Salesman methods
func (r *surveyRepositoryImpl) StoreSalesmen(tx *sqlx.Tx, surveySalesmen []model.SurveySalesman) error {
	for _, surveySalesman := range surveySalesmen {
		query := `INSERT INTO mst.m_survey_salesman (cust_id, survey_id, salesman_id) VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, surveySalesman.CustId, surveySalesman.SurveyId, surveySalesman.SalesmanId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *surveyRepositoryImpl) DeleteSalesmenBySurveyId(tx *sqlx.Tx, surveyId int) error {
	query := `UPDATE mst.m_survey_salesman SET is_del = true WHERE survey_id = $1`
	_, err := tx.Exec(query, surveyId)
	return err
}

func (r *surveyRepositoryImpl) FindSalesmenBySurveyId(surveyId int) ([]model.SurveySalesman, error) {
	var salesmen []model.SurveySalesman
	query := `SELECT ss.m_survey_salesman_id,
		ss.cust_id,
		ss.survey_id,
		ss.salesman_id,
		COALESCE(sm.sales_team_id, fallback_sm.sales_team_id) AS sales_team_id,
		st.sales_team_name,
		COALESCE(sm.sales_name, fallback_sm.sales_name) AS sales_name
		FROM mst.m_survey_salesman ss
		LEFT JOIN mst.m_salesman sm ON sm.emp_id = ss.salesman_id AND sm.cust_id = ss.cust_id AND sm.is_del = false
		LEFT JOIN LATERAL (
			SELECT fallback.emp_id, fallback.cust_id, fallback.sales_name, fallback.sales_team_id
			FROM mst.m_salesman fallback
			JOIN smc.m_customer mc ON mc.cust_id = fallback.cust_id
			JOIN mst.m_survey_area sa ON sa.survey_id = ss.survey_id
				AND sa.distributor_id = mc.distributor_id
				AND sa.is_del = false
			WHERE fallback.emp_id = ss.salesman_id
				AND fallback.is_del = false
			ORDER BY fallback.cust_id ASC
			LIMIT 1
		) fallback_sm ON sm.emp_id IS NULL
		LEFT JOIN mst.m_sales_team st ON st.sales_team_id = COALESCE(sm.sales_team_id, fallback_sm.sales_team_id)
			AND st.cust_id = COALESCE(sm.cust_id, fallback_sm.cust_id)
		WHERE ss.survey_id = $1 AND ss.is_del = false
		ORDER BY ss.m_survey_salesman_id ASC`
	err := r.DB.Select(&salesmen, query, surveyId)
	return salesmen, err
}

// Survey Distributor methods
func (r *surveyRepositoryImpl) StoreSurveyDistributors(tx *sqlx.Tx, distributorIds []int, surveyId int, custId string, createdBy int64) error {
	now := time.Now().UTC()
	for _, distributorId := range distributorIds {
		query := `INSERT INTO mst.m_survey_distributor (cust_id, survey_id, distributor_id, created_at, created_by, updated_at, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := tx.Exec(query, custId, surveyId, distributorId, now, createdBy, now, createdBy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *surveyRepositoryImpl) DeleteSurveyDistributorsBySurveyId(tx *sqlx.Tx, surveyId int) error {
	query := `UPDATE mst.m_survey_distributor SET is_del = true, deleted_at = NOW() WHERE survey_id = $1 AND is_del = false`
	_, err := tx.Exec(query, surveyId)
	return err
}

func (r *surveyRepositoryImpl) FindSurveyDistributorsBySurveyId(surveyId int) ([]model.SurveyDistributor, error) {
	var distributors []model.SurveyDistributor
	query := `SELECT sd.m_survey_distributor_id, sd.cust_id, sd.survey_id, sd.distributor_id, sd.is_del,
		md.distributor_code, md.distributor_name
		FROM mst.m_survey_distributor sd
		LEFT JOIN mst.m_distributor md ON md.distributor_id = sd.distributor_id AND md.is_del = false
		WHERE sd.survey_id = $1 AND sd.is_del = false
		ORDER BY sd.m_survey_distributor_id ASC`
	err := r.DB.Select(&distributors, query, surveyId)
	return distributors, err
}

// Survey Outlet methods
func (r *surveyRepositoryImpl) StoreOutlets(tx *sqlx.Tx, surveyId int, outletIds []int) error {
	for _, outletId := range outletIds {
		query := `INSERT INTO mst.m_survey_outlet (survey_id, outlet_id) VALUES ($1, $2)`
		_, err := tx.Exec(query, surveyId, outletId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *surveyRepositoryImpl) DeleteOutletsBySurveyId(tx *sqlx.Tx, surveyId int) error {
	query := `UPDATE mst.m_survey_outlet SET is_del = true WHERE survey_id = $1`
	_, err := tx.Exec(query, surveyId)
	return err
}

func (r *surveyRepositoryImpl) FindOutletsBySurveyId(surveyId int) ([]model.SurveyOutlet, error) {
	var outlets []model.SurveyOutlet
	query := `SELECT so.survey_outlet_id,
		so.survey_id,
		so.outlet_id,
		COALESCE(o.outlet_principal_code, o.outlet_code) AS outlet_code,
		o.outlet_name,
		o.ot_class_id,
		oc.ot_class_name,
		o.ot_grp_id,
		og.ot_grp_name,
		o.ot_type_id,
		ot.ot_type_name
		FROM mst.m_survey_outlet so
		JOIN mst.m_survey s ON s.survey_id = so.survey_id AND s.is_del = false
		LEFT JOIN mst.m_outlet o ON so.outlet_id = o.outlet_id AND o.is_del = false
		LEFT JOIN smc.m_customer mc ON mc.cust_id = o.cust_id
		LEFT JOIN mst.m_outlet_class oc ON oc.ot_class_id = o.ot_class_id AND oc.cust_id = COALESCE(NULLIF(mc.parent_cust_id, ''), s.cust_id)
		LEFT JOIN mst.m_outlet_group og ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = COALESCE(NULLIF(mc.parent_cust_id, ''), s.cust_id)
		LEFT JOIN mst.m_outlet_type ot ON ot.ot_type_id = o.ot_type_id AND ot.cust_id = COALESCE(NULLIF(mc.parent_cust_id, ''), s.cust_id)
		WHERE so.survey_id = $1 AND so.is_del = false
		ORDER BY so.survey_outlet_id ASC`
	err := r.DB.Select(&outlets, query, surveyId)
	return outlets, err
}

// Survey Detail methods
func (r *surveyRepositoryImpl) StoreDetails(tx *sqlx.Tx, surveyId int, templateIds []int) error {
	for _, templateId := range templateIds {
		query := `INSERT INTO mst.m_survey_detail (survey_id, survey_template_id) VALUES ($1, $2)`
		_, err := tx.Exec(query, surveyId, templateId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *surveyRepositoryImpl) DeleteDetailsBySurveyId(tx *sqlx.Tx, surveyId int) error {
	query := `UPDATE mst.m_survey_detail SET is_del = true WHERE survey_id = $1`
	_, err := tx.Exec(query, surveyId)
	return err
}

func (r *surveyRepositoryImpl) FindDetailsBySurveyId(surveyId int) ([]model.SurveyDetail, error) {
	var details []model.SurveyDetail
	query := `SELECT survey_detail_id, survey_id, survey_template_id 
		FROM mst.m_survey_detail WHERE survey_id = $1 AND is_del = false`
	err := r.DB.Select(&details, query, surveyId)
	return details, err
}
