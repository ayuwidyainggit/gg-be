package utils

import (
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func CheckTableNotExists(db *gorm.DB, modelInstance interface{}) {
	if !db.Migrator().HasTable(modelInstance) {
		db.AutoMigrate(modelInstance)
	}
}

func AutoMigrate(db *gorm.DB) {
	CheckTableNotExists(db, &model.Pjp{})
	CheckTableNotExists(db, &model.Route{})
	CheckTableNotExists(db, &model.RoutePopPermanent{})
	CheckTableNotExists(db, &model.RouteOutlet{})
	CheckTableNotExists(db, &model.RoutePopDaily{})
}
