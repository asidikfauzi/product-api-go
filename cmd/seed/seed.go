package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"product-api-go/internal/database"
	"product-api-go/internal/model"
	"time"
)

func seedProductMeasurements(db *gorm.DB) ([]model.ProductMeasurements, error) {
	measurements := []model.ProductMeasurements{
		{Name: "Kilogram"},
		{Name: "Gram"},
		{Name: "Milligram"},
		{Name: "Liter"},
		{Name: "Pieces"},
		{Name: "Dozen"},
		{Name: "Gross"},
		{Name: "Ball"},
		{Name: "Karton"},
	}

	for _, measurement := range measurements {
		if err := db.FirstOrCreate(&measurement, model.ProductMeasurements{Name: measurement.Name}).Error; err != nil {
			return nil, err
		}
	}

	return measurements, nil
}

func seedCategories(db *gorm.DB) ([]model.Categories, error) {
	categories := []model.Categories{
		{Name: "Beverages", CreatedAt: time.Now()},
		{Name: "Snacks", CreatedAt: time.Now()},
		{Name: "Groceries", CreatedAt: time.Now()},
		{Name: "Vegetables", CreatedAt: time.Now()},
		{Name: "Protein", CreatedAt: time.Now()},
		{Name: "Fruit", CreatedAt: time.Now()},
	}

	for _, category := range categories {
		if err := db.FirstOrCreate(&category, model.Categories{Name: category.Name}).Error; err != nil {
			return nil, err
		}
	}

	return categories, nil
}

func main() {
	db := database.InitDatabase()

	measurements, err := seedProductMeasurements(db)
	if err != nil {
		log.Fatalf("Failed to seed product measurements: %v", err)
	}

	jsonMeasurements, err := json.MarshalIndent(measurements, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonMeasurements))

	categories, err := seedCategories(db)
	if err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}

	jsonCategories, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonCategories))

	fmt.Println("Database seeding completed successfully!")
}
