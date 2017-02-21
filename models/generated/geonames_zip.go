package models

type GeonamesZip struct {
	CountryID       int              `gorm:"primary_key" json:"country_id"`
	ProvinceCode    string           `json:"province_code"`
	Latitude        float64          `json:"latitude"`
	Longitude       float64          `json:"longitude"`
	ProvinceName    string           `json:"province_name"`
	Accuracy        string           `json:"accuracy"`
	PlaceName       string           `gorm:"primary_key" json:"place_name"`
	PostalCode      string           `gorm:"primary_key" json:"postal_code"`
	StateName       string           `json:"state_name"` // This is test comment for the field
	StateCode       string           `json:"state_code"`
	GeonamesCountry *GeonamesCountry `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // Belongs to "GeonamesCountry" with relation One-To-One.
	GeonamesState   *GeonamesState   `gorm:"ForeignKey:state_id;AssociationForeignKey:id"`   // Belongs to "GeonamesState" with relation Many-To-One.
	StateID         int              `json:"state_id"`                                       // Value holder field for relation "GeonamesState".
	GeonamesCity    *GeonamesCity    `gorm:"ForeignKey:city_id;AssociationForeignKey:id"`    // Belongs to "GeonamesCity" with relation Many-To-One.
	CityID          int              `json:"city_id"`                                        // Value holder field for relation "GeonamesCity".

}

// set GeonamesZip's table name to be `geonames_zips`
func (GeonamesZip) TableName() string {
	return "geonames_zips"
}
