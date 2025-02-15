package dto

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/pkg/response"
)

type ProductsResponseWithPage struct {
	Data []ProductResponse           `json:"data"`
	Page response.PaginationResponse `json:"page"`
}

type ProductResponse struct {
	ID            uuid.UUID              `json:"id"`
	Name          string                 `json:"name"`
	Description   *string                `json:"description,omitempty"`
	PurchasePrice int                    `json:"purchase_price"`
	SellingPrice  int                    `json:"selling_price"`
	TotalStock    int                    `json:"total_stock"`
	MinimumStock  int                    `json:"minimum_stock"`
	Image         *string                `json:"image"`
	Measurement   string                 `json:"measurement,omitempty"`
	Categories    []dto.CategoryResponse `json:"categories,omitempty"`
	CreatedAt     *string                `json:"created_at,omitempty"`
	CreatedBy     *uuid.UUID             `json:"created_by,omitempty"`
	UpdatedAt     *string                `json:"updated_at,omitempty"`
	UpdatedBy     *uuid.UUID             `json:"updated_by,omitempty"`
}
