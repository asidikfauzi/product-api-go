package main

import (
	"fmt"
	"product-api-go/internal/database"
)

func main() {
	database.InitDatabase()
	fmt.Println("Hello World")
}
