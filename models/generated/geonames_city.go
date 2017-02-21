package models

type GeonamesCity struct {
	ID          int     `gorm:"primary_key" json:"id"`
	Name        string  `json:"name"`
	Asciiname   string  `json:"asciiname"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	FeatureCode string  `json:"feature_code"`
	Population  int64   `json:"population"`
}

// set GeonamesCity's table name to be `geonames_cities`
func (GeonamesCity) TableName() string {
	return "geonames_cities"
}
