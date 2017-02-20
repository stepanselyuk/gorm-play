package models

import "time"

type Product struct {
	ID        int        `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Code      string     `json:"code"`
	Price     uint       `json:"price"`
	Vendor    string     `json:"vendor"`
}

// set User's table name to be `profiles`
func (Product) TableName() string {
	return "profiles"
}