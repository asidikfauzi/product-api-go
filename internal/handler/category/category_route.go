package category

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, categories *CategoriesController) {
	categoryGroup := r.Group("/categories")
	categoryGroup.GET("", categories.FindAll)
	categoryGroup.GET("/:id", categories.FindById)
}
