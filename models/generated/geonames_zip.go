package models

type GeonamesZip struct {
	CountryID    int     `gorm:"primary_key" json:"country_id"`
	PlaceName    string  `gorm:"primary_key" json:"place_name"`
	PostalCode   string  `gorm:"primary_key" json:"postal_code"`
	StateName    string  `json:"state_name"`
	StateCode    string  `json:"state_code"`
	ProvinceName string  `json:"province_name"`
	ProvinceCode string  `json:"province_code"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Accuracy     string  `json:"accuracy"`
}

// set GeonamesZip's table name to be `geonames_zips`
func (GeonamesZip) TableName() string {
	return "geonames_zips"
}
