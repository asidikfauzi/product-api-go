package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"product-api-go/internal/handler/category"
	"product-api-go/internal/handler/measurement"
	"product-api-go/internal/handler/product"
	"product-api-go/internal/injector"
)

type Server struct {
	Engine *gin.Engine
}

func InitializedServer() *Server {
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		fmt.Println("pong")
		c.String(200, "pong")
	})

	categoryModule := injector.InitializedCategoriesModule()
	category.RegisterRoutes(api, categoryModule)

	measurementModule := injector.InitializedMeasurementsModule()
	measurement.RegisterRoutes(api, measurementModule)

	productModule := injector.InitializedProductsModule()
	product.RegisterRoutes(api, productModule)

	return &Server{Engine: r}
}
