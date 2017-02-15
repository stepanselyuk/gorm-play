package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {

	db, err := gorm.Open("mysql", "gorm-play:1234@tcp(gorm-db:3306)/gorm-play?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic(fmt.Sprint(err))
	}

	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1200", Price: 1000})
	db.Create(&Product{Code: "L1201", Price: 2000})
	db.Create(&Product{Code: "L1202", Price: 3000})

	// Read
	var product Product

	//db.First(&product, 4) // find product with id 1
	db.First(&product, "code = ?", "L1200") // find product with code l1212

	// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 3000)

	fmt.Println(product)

	// Delete - delete product
	//db.Delete(&product)
}
