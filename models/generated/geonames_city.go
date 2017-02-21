package models

type GeonamesCity struct {
	ID              int              `gorm:"primary_key" json:"id"`
	Name            string           `json:"name"`
	Asciiname       string           `json:"asciiname"`
	Latitude        float64          `json:"latitude"`
	Longitude       float64          `json:"longitude"`
	FeatureCode     string           `json:"feature_code"`
	Population      int64            `json:"population"`
	GeonamesState   *GeonamesState   `gorm:"ForeignKey:state_id;AssociationForeignKey:id"`   // Belongs to "GeonamesState" with relation Many-To-One.
	StateID         int              `json:"state_id"`                                       // Value holder field for relation "GeonamesState".
	GeonamesCountry *GeonamesCountry `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // Belongs to "GeonamesCountry" with relation Many-To-One.
	CountryID       int              `json:"country_id"`                                     // Value holder field for relation "GeonamesCountry".

}

// set GeonamesCity's table name to be `geonames_cities`
func (GeonamesCity) TableName() string {
	return "geonames_cities"
}
