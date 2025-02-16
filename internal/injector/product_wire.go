//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"product-api-go/internal/database"
	handler "product-api-go/internal/handler/product"
	postgresCategory "product-api-go/internal/repository/postgres/category"
	postgresMeasurement "product-api-go/internal/repository/postgres/measurement"
	postgresProduct "product-api-go/internal/repository/postgres/product"
	redisProduct "product-api-go/internal/repository/redis/product"
)

func InitializedProductsModule() *handler.ProductsController {
	wire.Build(
		database.InitDatabase,
		database.InitRedis,
		postgresProduct.NewProductsPostgres,
		redisProduct.NewProductsRedis,
		postgresMeasurement.NewMeasurementsPostgres,
		postgresCategory.NewCategoriesPostgres,
		handler.NewProductsService,
		handler.NewProductsController,
	)

	return nil
}
