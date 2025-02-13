//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"product-api-go/internal/database"
	handler "product-api-go/internal/handler/category"
	postgres "product-api-go/internal/repository/postgres/category"
)

func InitializedCategoriesModule() *handler.CategoriesController {
	wire.Build(
		database.InitDatabase,
		postgres.NewCategoriesPostgres,
		handler.NewCategoriesService,
		handler.NewCategoriesController,
	)

	return nil
}
