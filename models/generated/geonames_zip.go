package models_generated

type GeonamesZip struct {
	ProvinceCode    string           `json:"province_code"`
	Accuracy        string           `json:"accuracy"`
	PlaceName       string           `gorm:"primary_key" json:"place_name"`
	ProvinceName    string           `json:"province_name"`
	StateName       string           `json:"state_name"` // This is test comment for the field
	StateCode       string           `json:"state_code"`
	Latitude        float64          `json:"latitude"`
	Longitude       float64          `json:"longitude"`
	CountryID       int              `gorm:"primary_key" json:"country_id"`
	PostalCode      string           `gorm:"primary_key" json:"postal_code"`
	GeonamesCountry *GeonamesCountry `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // BelongsTo "GeonamesCountry" with relation Many(GeonamesZip)-To-One(GeonamesCountry).
	GeonamesState   *GeonamesState   `gorm:"ForeignKey:state_id;AssociationForeignKey:id"`   // BelongsTo "GeonamesState" with relation Many(GeonamesZip)-To-One(GeonamesState).
	StateID         int              `json:"state_id"`                                       // Value holder field for relation with "GeonamesState".
	GeonamesCity    *GeonamesCity    `gorm:"ForeignKey:city_id;AssociationForeignKey:id"`    // BelongsTo "GeonamesCity" with relation Many(GeonamesZip)-To-One(GeonamesCity).
	CityID          int              `json:"city_id"`                                        // Value holder field for relation with "GeonamesCity".

}

// set GeonamesZip's table name to be `geonames_zips`
func (GeonamesZip) TableName() string {
	return "geonames_zips"
}
