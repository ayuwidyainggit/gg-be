package outlet_visit_principle

import (
	"context"
	"gorm.io/gorm"
	"scyllax-pjp/model"
)

func (repo *outletVisitPrincipleRepository) UpdateOutletVisitListWithFile(ctx context.Context, tx *gorm.DB, id int64, date string, fileInfo model.OutletVisitListPrinciple) error {
	result := tx.WithContext(ctx).Exec(`
		UPDATE pjp_principles.outlet_visit_list
		SET 
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
			location_status = ?
		WHERE date = ? AND id = ?
	`,
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
		date,
		id,
	)

	return result.Error
}
