package main

import (
	"fmt"
	"product-api-go/internal/config"
	"product-api-go/internal/database"
	"product-api-go/internal/server"
)

func main() {
	fmt.Println("Starting server...")
	database.InitDatabase()

	s := server.InitializedServer()

	err := s.Engine.Run(fmt.Sprintf(":%s", config.Env("APP_PORT")))
	if err != nil {
		return
	}
}
