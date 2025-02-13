package dto

import (
	"github.com/google/uuid"
	"product-api-go/internal/pkg/response"
)

type CategoriesQuery struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Search    string `form:"search"`
	OrderBy   string `form:"order_by"`
	Direction string `form:"direction"`
}

type FindAllCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type FindByIdCategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt string    `json:"updated_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
}

type FindAllCategoriesResponse struct {
	Data []FindAllCategory           `json:"data"`
	Page response.PaginationResponse `json:"page"`
}
