package models_generated

type GeonamesCountry struct {
	Name           string           `json:"name"`
	Code           string           `json:"code"`
	Population     int64            `json:"population"`
	ID             int              `gorm:"primary_key" json:"id"`
	GeonamesStates []*GeonamesState `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // HasMany "GeonamesState" with relation One(GeonamesCountry)-Has-Many(GeonamesState).
	GeonamesZips   []*GeonamesZip   `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // HasMany "GeonamesZip" with relation One(GeonamesCountry)-Has-Many(GeonamesZip).
	GeonamesCities []*GeonamesCity  `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // HasMany "GeonamesCity" with relation One(GeonamesCountry)-Has-Many(GeonamesCity).

}

// set GeonamesCountry's table name to be `geonames_countries`
func (GeonamesCountry) TableName() string {
	return "geonames_countries"
}
