//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"product-api-go/internal/database"
	handler "product-api-go/internal/handler/category"
	postgres "product-api-go/internal/repository/postgres/category"
	redis "product-api-go/internal/repository/redis/category"
)

func InitializedCategoriesModule() *handler.CategoriesController {
	wire.Build(
		database.InitDatabase,
		database.InitRedis,
		postgres.NewCategoriesPostgres,
		redis.NewCategoriesRedis,
		handler.NewCategoriesService,
		handler.NewCategoriesController,
	)

	return nil
}
