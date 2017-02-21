package models

type GeonamesState struct {
	ID        int    `gorm:"primary_key" json:"id"`
	Name      string `json:"name"`
	Asciiname string `json:"asciiname"`
	Code      string `json:"code"`
}

// set GeonamesState's table name to be `geonames_states`
func (GeonamesState) TableName() string {
	return "geonames_states"
}
