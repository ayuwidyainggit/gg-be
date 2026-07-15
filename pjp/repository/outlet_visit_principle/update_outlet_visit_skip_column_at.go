package outlet_visit_principle

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *outletVisitPrincipleRepository) UpdateOutletVisitListSkipColumnAt(ctx context.Context, tx *gorm.DB, column string, currentTime int64, date string, id int64, skipReson string, inOutlet bool, fileInfo model.OutletVisitListPrinciple) {
	if fileInfo.FileUrl == "" {
		result := tx.WithContext(ctx).Exec(`
		UPDATE pjp_principles.outlet_visit_list
		SET `+column+` = ?, skip_reason = ?, skip_in_outlet = ?
		WHERE date = ? AND id = ?
		`, currentTime, skipReson, inOutlet, date, id)
		helper.ErrorPanic(result.Error)
		return
	}

	result := tx.WithContext(ctx).Exec(`
		UPDATE pjp_principles.outlet_visit_list
		SET 
		    `+column+` = ?,
			skip_reason = ?,
			is_update_location = ?,
			file_name = ?,
			file_type = ?,
			media_category = ?,
			file_url = ?,
			file_size = ?,
			file_base64 = ?,
			photo_path = ?,
			folder = ?,
			latitude = ?,
			longitude = ?,
			allowed_radius = ?,
			distance_meter = ?,
			location_status = ?,
			skip_in_outlet = ?
		WHERE date = ? AND id = ?
	`,
		currentTime,
		skipReson,
		fileInfo.IsUpdateLocation,
		fileInfo.FileName,
		fileInfo.FileType,
		fileInfo.MediaCategory,
		fileInfo.FileUrl,
		fileInfo.FileSize,
		fileInfo.FileBase64,
		fileInfo.PhotoPath,
		fileInfo.Folder,
		fileInfo.Latitude,
		fileInfo.Longitude,
		fileInfo.AllowedRadius,
		fileInfo.DistanceMeter,
		fileInfo.LocationStatus,
		fileInfo.SkipInOutlet,
		date,
		id,
	)

	helper.ErrorPanic(result.Error)
}
