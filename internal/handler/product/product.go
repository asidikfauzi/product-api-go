package product

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/product/dto"
)

type ProductsService interface {
	FindAll(query dto.ProductQuery) (dto.ProductsResponseWithPage, int, error)
	FindById(id uuid.UUID) (dto.ProductResponse, int, error)
	Create(input dto.ProductInput) (dto.ProductResponse, int, error)
	Update(id uuid.UUID, input dto.ProductInput) (dto.ProductResponse, int, error)
	Delete(id uuid.UUID) (dto.ProductResponse, int, error)
}
