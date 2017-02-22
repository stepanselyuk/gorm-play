package models

import genmodels "github.com/stepanselyuk/gorm-play/models/generated"

type GeonamesZip struct {
	genmodels.GeonamesZip `gorm:"embedded"`
}
