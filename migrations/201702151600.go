package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/stepanselyuk/gorm-play/models"
	"gopkg.in/gormigrate.v1"
)

func migration201702151600() *gormigrate.Migration {

	// added Vendor field to models.Product

	return &gormigrate.Migration{
		ID: "201702151600",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Product{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropColumn("vendor").Error
		},
	}
}
