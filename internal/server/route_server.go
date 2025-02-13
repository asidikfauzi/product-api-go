package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"product-api-go/internal/handler/category"
	"product-api-go/internal/injector"
)

type Server struct {
	Engine *gin.Engine
}

func InitializedServer() *Server {
	r := gin.Default()

	api := r.Group("/api") // Buat grup dengan prefix "/api"

	api.GET("/ping", func(c *gin.Context) { // Sekarang "/api/ping"
		fmt.Println("pong")
		c.String(200, "pong")
	})

	categoryModule := injector.InitializedCategoriesModule()
	category.RegisterRoutes(api, categoryModule)

	return &Server{Engine: r}
}
