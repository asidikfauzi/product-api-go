package measurement

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, measurements *MeasurementsController) {
	measurementGroup := r.Group("/measurements")
	measurementGroup.GET("", measurements.FindAll)
	measurementGroup.GET("/:id", measurements.FindById)
	measurementGroup.POST("", measurements.Create)
	measurementGroup.PUT("/:id", measurements.Update)
}
