package models_generated

type GeonamesCity struct {
	FeatureCode     string           `json:"feature_code"`
	Population      int64            `json:"population"`
	ID              int              `gorm:"primary_key" json:"id"`
	Name            string           `json:"name"`
	Asciiname       string           `json:"asciiname"`
	Latitude        float64          `json:"latitude"`
	Longitude       float64          `json:"longitude"`
	GeonamesCountry *GeonamesCountry `gorm:"ForeignKey:country_id;AssociationForeignKey:id"` // BelongsTo "GeonamesCountry" with relation Many(GeonamesCity)-To-One(GeonamesCountry).
	CountryID       int              `json:"country_id"`                                     // Value holder field for relation with "GeonamesCountry".
	GeonamesZips    []*GeonamesZip   `gorm:"ForeignKey:city_id;AssociationForeignKey:id"`    // HasMany "GeonamesZip" with relation One(GeonamesCity)-Has-Many(GeonamesZip).
	GeonamesState   *GeonamesState   `gorm:"ForeignKey:state_id;AssociationForeignKey:id"`   // BelongsTo "GeonamesState" with relation Many(GeonamesCity)-To-One(GeonamesState).
	StateID         int              `json:"state_id"`                                       // Value holder field for relation with "GeonamesState".

}

// set GeonamesCity's table name to be `geonames_cities`
func (GeonamesCity) TableName() string {
	return "geonames_cities"
}
