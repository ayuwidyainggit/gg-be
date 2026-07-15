package live_monitoring

import (
	"context"
	"errors"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"time"

	"gorm.io/gorm"
)

func (s *liveMonitoringService) GetUpdateLocations(ctx context.Context, req request.UpdateLocationsRequest, custID string) (response.UpdateLocationsResponse, error) {
	date := req.Date
	if date == "" {
		date = time.Now().In(jakartaLocation).Format("2006-01-02")
	}
	resolvedCustID, err := s.repository.GetEmployeeRole(ctx, s.db, req.EmpID, custID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.UpdateLocationsResponse{}, gorm.ErrRecordNotFound
		}
		return response.UpdateLocationsResponse{}, err
	}
	branch := "pjp_principles"
	if len(resolvedCustID) > 6 {
		branch = "pjp"
	}
	rows, err := s.repository.GetUpdateLocations(ctx, s.db, req.EmpID, date, custID, branch)
	if err != nil {
		return response.UpdateLocationsResponse{}, err
	}
	timeline := make([]response.TimelineItem, 0, len(rows))
	for i, row := range rows {
		timeline = append(timeline, response.TimelineItem{
			Sequence: i + 1, Type: row.Type, Latitude: row.Latitude, Longitude: row.Longitude,
			DestinationID: row.DestinationID, DestinationType: row.DestinationType, DestinationName: row.DestinationName,
			RecordedAt: row.RecordedAt,
		})
	}
	return response.UpdateLocationsResponse{Timeline: timeline}, nil
}
