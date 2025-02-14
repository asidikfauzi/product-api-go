package dto

import (
	"github.com/google/uuid"
	"product-api-go/internal/pkg/response"
)

type MeasurementsResponseWithPage struct {
	Data []MeasurementResponse       `json:"data"`
	Page response.PaginationResponse `json:"page"`
}

type MeasurementResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
