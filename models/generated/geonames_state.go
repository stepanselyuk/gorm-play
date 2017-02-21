package models

type GeonamesState struct {
	ID              int              `gorm:"primary_key" json:"id"`
	Name            string           `json:"name"`
	Asciiname       string           `json:"asciiname"`
	Code            string           `json:"code"`
	GeonamesCountry *GeonamesCountry `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // Belongs to "GeonamesCountry" with relation Many-To-One.
	CountryID       int              `json:"country_id"`                                     // Value holder field for relation "GeonamesCountry".

}

// set GeonamesState's table name to be `geonames_states`
func (GeonamesState) TableName() string {
	return "geonames_states"
}
