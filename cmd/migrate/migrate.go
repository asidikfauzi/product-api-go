package main

import (
	"fmt"
	"product-api-go/internal/database"
	"product-api-go/internal/model"
)

const migrationErrorMsg = "Error during '%s' table migration"

func main() {
	db := database.InitDatabase()

	db.Debug().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := db.Debug().AutoMigrate(&model.ProductMeasurements{}); err != nil {
		panic(fmt.Sprintf(migrationErrorMsg, "ProductMeasurements"))
	}

	if err := db.Debug().AutoMigrate(&model.Categories{}); err != nil {
		panic(fmt.Sprintf(migrationErrorMsg, "Categories"))
	}

	if err := db.Debug().AutoMigrate(&model.Products{}); err != nil {
		panic(fmt.Sprintf(migrationErrorMsg, "Products"))
	}

	fmt.Println("Database migration successful")
}
