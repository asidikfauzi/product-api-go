package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

func Env(key string) string {

	appConfig, err := godotenv.Read()

	if err != nil {
		fmt.Println("Error reading .env file")
	}

	value, ok := appConfig[key]

	if !ok {
		fmt.Printf("Environment variable %s not found\n", key)
		return appConfig[key]
	}

	return value
}
