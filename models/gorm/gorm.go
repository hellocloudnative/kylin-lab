package gorm

import (
	"github.com/jinzhu/gorm"
	"kylin-lab/models"
)

func AutoMigrate(db *gorm.DB) error {
	db.SingularTable(true)
	return db.AutoMigrate(
		new(models.LabUser),
		new(models.LabVirtualMachine),
	).Error
}
