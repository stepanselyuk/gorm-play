package models

import genmodels "github.com/stepanselyuk/gorm-play/models/generated"

type GeonamesCountry struct {
	genmodels.GeonamesCountry `gorm:"embedded"`
}
