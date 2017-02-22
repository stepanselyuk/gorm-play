package models_usage

import (
	"fmt"

	"github.com/stepanselyuk/gorm-play/database"
	models "github.com/stepanselyuk/gorm-play/models"
)

//Called from cmd. Example of how we can structure things.
func FindCity() {

	db := database.GetConnection()
	defer db.Close()

	// db.Create(&models.Product{Code: "L1200", Price: 1000})
	//database.Model(&product).Update("Price", 3000)
	//database.Delete(&product)

	var city models.GeonamesCity

	db.Preload("GeonamesCountry").First(&city, 2759794) // find city with id
	//db.First(&product, "code = ?", "L1200") // find product with code l1212

	fmt.Println(fmt.Sprintf("City name: %p", city.Name))
	fmt.Println(fmt.Sprintf("Country name: %p", city.GeonamesCountry.Name))

	var country models.GeonamesCountry

	db.Preload("GeonamesStates").First(&country, "code = ?", "NL")

	fmt.Println(fmt.Sprintf("Country name: %p", country.Name))
	fmt.Println(fmt.Sprintf("Count of states: %p", len(country.GeonamesStates)))

	//fmt.Println(city.GeonamesCountry.Name)
}
