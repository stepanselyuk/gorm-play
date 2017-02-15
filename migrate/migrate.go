package migrate

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/stepanselyuk/gorm-play/database"
	"github.com/stepanselyuk/gorm-play/models"
	"gopkg.in/gormigrate.v1"
)

//Called from cmd. Example of how we can structure things.
func Example() {
	//fmt.Println("Hi there. Would this work :)?")
	migration()
}

func migration() {

	db := database.GetConnection()

	defer db.Close()

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// added Vendor field to models.Product
		{
			ID: "201702151600",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Product{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropColumn("vendor").Error
			},
		},
	})

	err := m.Migrate()
	if err == nil {
		log.Printf("Migration did run successfully")
	} else {
		log.Printf("Could not migrate: %v", err)
	}
}
