package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var connection *gorm.DB

func GetConnection() *gorm.DB {

	if connection != nil {
		pingDb(connection)
		return connection
	}

	connection, err := gorm.Open("mysql", "gorm-play:1234@tcp(172.19.0.2:3306)/gorm-play?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatal(err)
	}

	pingDb(connection)
	connection.LogMode(true)

	return connection
}

func pingDb(db *gorm.DB) {

	if err := db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
}
