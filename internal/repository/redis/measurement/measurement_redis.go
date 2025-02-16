package measurement

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/measurement/dto"
)

type MeasurementsRedis interface {
	GetAllMeasurement(key string) (dto.MeasurementsResponseWithPage, error)
	CreateAllMeasurement(key string, data dto.MeasurementsResponseWithPage) error
	GetMeasurementById(key uuid.UUID) (dto.MeasurementResponse, error)
	CreateMeasurementById(key uuid.UUID, data dto.MeasurementResponse) error
	DeleteAll(key string) error
	Delete(key string) error
}
