package measurement

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/measurement/dto"
	"product-api-go/internal/model"
)

type MeasurementsPostgres interface {
	FindAll(q dto.MeasurementQuery) ([]model.ProductMeasurements, int64, error)
	FindById(id uuid.UUID) (model.ProductMeasurements, error)
	FindByName(name string) (model.ProductMeasurements, error)
	Create(input dto.MeasurementInput) (model.ProductMeasurements, error)
	Update(id uuid.UUID, input dto.MeasurementInput) (model.ProductMeasurements, error)
}
