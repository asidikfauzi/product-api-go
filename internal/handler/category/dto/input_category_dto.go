package dto

type CategoryInput struct {
	Name string `json:"name" validate:"required"`
}
