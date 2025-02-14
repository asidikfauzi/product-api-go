package category

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, categories *CategoriesController) {
	categoryGroup := r.Group("/categories")
	categoryGroup.GET("", categories.FindAll)
	categoryGroup.GET("/:id", categories.FindById)
	categoryGroup.POST("", categories.Create)
	categoryGroup.PUT("/:id", categories.Update)
	categoryGroup.DELETE("/:id", categories.Delete)
}
