//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"product-api-go/internal/database"
	handler "product-api-go/internal/handler/measurement"
	postgres "product-api-go/internal/repository/postgres/measurement"
)

func InitializedMeasurementsModule() *handler.MeasurementsController {
	wire.Build(
		database.InitDatabase,
		postgres.NewMeasurementsPostgres,
		handler.NewMeasurementsService,
		handler.NewMeasurementsController,
	)

	return nil
}
