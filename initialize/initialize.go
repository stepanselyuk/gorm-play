package initialize

import (
	"fmt"

	"github.com/stepanselyuk/gorm-play/database"
	"github.com/stepanselyuk/gorm-play/models"
)

//Called from cmd. Example of how we can structure things.
func Example() {
	//fmt.Println("Hi there. Would this work :)?")
	createProducts()
}

func createProducts() {

	db := database.GetConnection()
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&models.Product{})

	// Create
	db.Create(&models.Product{Code: "L1200", Price: 1000})
	db.Create(&models.Product{Code: "L1201", Price: 2000})
	db.Create(&models.Product{Code: "L1202", Price: 3000})

	return

	// Read
	var product models.Product

	//database.First(&product, 4) // find product with id 1
	db.First(&product, "code = ?", "L1200") // find product with code l1212

	// Update - update product's price to 2000
	//database.Model(&product).Update("Price", 3000)

	fmt.Println(product)

	// Delete - delete product
	//database.Delete(&product)
}
