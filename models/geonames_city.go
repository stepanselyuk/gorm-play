package models

import genmodels "github.com/stepanselyuk/gorm-play/models/generated"

type GeonamesCity struct {
	genmodels.GeonamesCity `gorm:"embedded"`
}
