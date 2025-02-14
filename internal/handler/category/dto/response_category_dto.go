package dto

import (
	"github.com/google/uuid"
	"product-api-go/internal/pkg/response"
)

type CategoriesResponseWithPage struct {
	Data []CategoryResponse          `json:"data"`
	Page response.PaginationResponse `json:"page"`
}

type CategoryResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *string    `json:"created_at,omitempty"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt *string    `json:"updated_at,omitempty"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}
