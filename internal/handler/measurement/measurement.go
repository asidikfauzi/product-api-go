package measurement

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/measurement/dto"
)

type MeasurementsService interface {
	FindAll(query dto.MeasurementQuery) (dto.MeasurementsResponseWithPage, int, error)
	FindById(id uuid.UUID) (dto.MeasurementResponse, int, error)
	Create(input dto.MeasurementInput) (dto.MeasurementResponse, int, error)
	Update(id uuid.UUID, input dto.MeasurementInput) (dto.MeasurementResponse, int, error)
}
