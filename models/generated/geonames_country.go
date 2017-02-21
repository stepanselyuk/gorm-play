package models

type GeonamesCountry struct {
	ID         int    `gorm:"primary_key" json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Population int64  `json:"population"`
}

// set GeonamesCountry's table name to be `geonames_countries`
func (GeonamesCountry) TableName() string {
	return "geonames_countries"
}
