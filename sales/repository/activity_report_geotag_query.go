package repository

import (
	"fmt"
	"time"

	"sales/model"

	"gorm.io/gorm"
)

func buildActivityReportGeotagSQL(schema string, withEmpFilter bool) string {
	empFilter := ""
	if withEmpFilter {
		empFilter = "\n        AND p.salesman_id = ?"
	}

	return fmt.Sprintf(`
WITH pjp_data AS (
    SELECT
        v.location_status,
        p.salesman_id,
        p.salesman_name
    FROM %s.outlet_visit_list v
    JOIN %s.permanent_journey_plans p ON v.pjp_code = p.pjp_code
    WHERE
        p.cust_id IN ?
        AND v.date::date BETWEEN ?::date AND ?::date%s
)
SELECT
    p.salesman_id AS salesman_code,
    p.salesman_name,
    COUNT(*) AS total_visit,
    SUM(CASE WHEN p.location_status = 1 THEN 1 ELSE 0 END) AS geotag_match_count,
    SUM(CASE WHEN p.location_status = 0 THEN 1 ELSE 0 END) AS geotag_unmatch_count,
    ROUND(
        100.0 * SUM(CASE WHEN p.location_status = 1 THEN 1 ELSE 0 END)
        / NULLIF(COUNT(*), 0),
        2
    ) AS geotag_match_pct,
    ROUND(
        100.0 * SUM(CASE WHEN p.location_status = 0 THEN 1 ELSE 0 END)
        / NULLIF(COUNT(*), 0),
        2
    ) AS geotag_unmatch_pct
FROM pjp_data p
GROUP BY
    p.salesman_id,
    p.salesman_name
ORDER BY
    p.salesman_name`, schema, schema, empFilter)
}

func activityReportGeotagDateRange(year int) (time.Time, time.Time) {
	dateStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)
	return dateStart, dateEnd
}

func queryActivityReportGeotag(
	db *gorm.DB,
	schema string,
	custIDs []string,
	year int,
	empID *int,
) ([]model.ActivityReportGeotagRow, error) {
	if len(custIDs) == 0 {
		return nil, nil
	}

	dateStart, dateEnd := activityReportGeotagDateRange(year)
	withEmpFilter := empID != nil
	args := []interface{}{custIDs, dateStart, dateEnd}
	if withEmpFilter {
		args = append(args, *empID)
	}

	var rows []model.ActivityReportGeotagRow
	err := db.
		Raw(
			buildActivityReportGeotagSQL(schema, withEmpFilter),
			args...,
		).
		Find(&rows).Error
	return rows, err
}

func mergeActivityReportGeotagRows(rows []model.ActivityReportGeotagRow) []model.ActivityReportGeotagRow {
	if len(rows) == 0 {
		return rows
	}

	type aggKey struct {
		salesmanCode int64
		salesmanName string
	}

	merged := make(map[aggKey]*model.ActivityReportGeotagRow)
	order := make([]aggKey, 0, len(rows))

	for _, row := range rows {
		key := aggKey{salesmanCode: row.SalesmanCode, salesmanName: row.SalesmanName}
		existing, ok := merged[key]
		if !ok {
			copied := row
			merged[key] = &copied
			order = append(order, key)
			continue
		}
		existing.TotalVisit += row.TotalVisit
		existing.GeotagMatchCount += row.GeotagMatchCount
		existing.GeotagUnmatchCount += row.GeotagUnmatchCount
	}

	results := make([]model.ActivityReportGeotagRow, 0, len(order))
	for _, key := range order {
		row := merged[key]
		if row.TotalVisit > 0 {
			row.GeotagMatchPct = roundActivityReportGeotagPct(float64(row.GeotagMatchCount), float64(row.TotalVisit))
			row.GeotagUnmatchPct = roundActivityReportGeotagPct(float64(row.GeotagUnmatchCount), float64(row.TotalVisit))
		}
		results = append(results, *row)
	}

	return results
}

func roundActivityReportGeotagPct(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return float64(int64((100.0*part/total+0.005)*100)) / 100
}

func (repository *RepositoryReportImpl) ActivityReportGeotag(
	parentCustID string,
	custIDs []string,
	year int,
	empID *int,
) ([]model.ActivityReportGeotagRow, error) {
	return activityReportGeotag(repository, parentCustID, custIDs, year, empID)
}

func activityReportGeotag(
	repository *RepositoryReportImpl,
	parentCustID string,
	custIDs []string,
	year int,
	empID *int,
) ([]model.ActivityReportGeotagRow, error) {
	principalIDs, distributorIDs := splitActivityReportCustIDs(parentCustID, custIDs)

	var principalScope, distributorScope []string
	if len(principalIDs) > 0 && hasActivityReportPrincipalPJP(repository.DB, principalIDs) {
		principalScope = principalIDs
	} else {
		distributorScope = append(distributorScope, principalIDs...)
	}
	distributorScope = append(distributorScope, distributorIDs...)

	var allRows []model.ActivityReportGeotagRow

	if len(principalScope) > 0 {
		rows, err := queryActivityReportGeotag(repository.DB, "pjp_principles", principalScope, year, empID)
		if err != nil {
			return nil, err
		}
		allRows = append(allRows, rows...)
	}

	if len(distributorScope) > 0 {
		rows, err := queryActivityReportGeotag(repository.DB, "pjp", distributorScope, year, empID)
		if err != nil {
			return nil, err
		}
		allRows = append(allRows, rows...)
	}

	return mergeActivityReportGeotagRows(allRows), nil
}