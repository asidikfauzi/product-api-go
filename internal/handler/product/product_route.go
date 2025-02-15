package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, products *ProductsController) {
	productGroup := r.Group("/products")
	productGroup.GET("", products.FindAll)
	productGroup.GET("/:id", products.FindById)
	productGroup.POST("", products.Create)
	productGroup.PUT("/:id", products.Update)
	productGroup.DELETE("/:id", products.Delete)
}
