package dto

type MeasurementInput struct {
	Name string `json:"name" validate:"required"`
}
