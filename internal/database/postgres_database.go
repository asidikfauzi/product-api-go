package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"product-api-go/internal/config"
	"time"
)

func InitDatabase() *gorm.DB {
	dbConfig := config.LoadDBConfigFromEnv()
	dbName := dbConfig.DBName
	dbUser := dbConfig.User
	dbPass := dbConfig.Password
	dbHost := dbConfig.Host
	dbPort := dbConfig.Port
	dbSSLMode := dbConfig.SSLMode
	dbTimeZone := dbConfig.TimeZone

	initDSN := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPass, dbPort, dbSSLMode, dbTimeZone,
	)

	db, err := gorm.Open(postgres.Open(initDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	var exists int
	checkDBQuery := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbName)
	db.Raw(checkDBQuery).Scan(&exists)

	if exists == 0 {
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
		if err := db.Exec(createDBQuery).Error; err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Println("Database created successfully:", dbName)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	finalDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSSLMode, dbTimeZone,
	)

	db, err = gorm.Open(postgres.Open(finalDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to newly created database: %v", err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		log.Fatalf("Failed to access database connection pool: %v", err)
	}

	sqlDB.SetConnMaxIdleTime(10 * time.Second)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	fmt.Println("Database connection established successfully!")
	return db
}
