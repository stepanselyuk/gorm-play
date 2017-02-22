package models_generated

type GeonamesState struct {
	Asciiname       string           `json:"asciiname"`
	Code            string           `json:"code"`
	ID              int              `gorm:"primary_key" json:"id"`
	Name            string           `json:"name"`
	GeonamesCountry *GeonamesCountry `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // BelongsTo "GeonamesCountry" with relation Many(GeonamesState)-To-One(GeonamesCountry).
	CountryID       int              `json:"country_id"`                                     // Value holder field for relation with "GeonamesCountry".
	GeonamesZips    []*GeonamesZip   `gorm:"ForeignKey:state_id;AssociationForeignKey:id"`   // HasMany "GeonamesZip" with relation One(GeonamesState)-Has-Many(GeonamesZip).
	GeonamesCities  []*GeonamesCity  `gorm:"ForeignKey:state_id;AssociationForeignKey:id"`   // HasMany "GeonamesCity" with relation One(GeonamesState)-Has-Many(GeonamesCity).

}

// set GeonamesState's table name to be `geonames_states`
func (GeonamesState) TableName() string {
	return "geonames_states"
}
