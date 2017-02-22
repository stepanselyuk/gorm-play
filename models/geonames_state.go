package models

import genmodels "github.com/stepanselyuk/gorm-play/models/generated"

type GeonamesState struct {
	genmodels.GeonamesState `gorm:"embedded"`
}
