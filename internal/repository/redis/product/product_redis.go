package product

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/product/dto"
)

type ProductsRedis interface {
	GetAllProduct(key string) (dto.ProductsResponseWithPage, error)
	CreateAllProduct(key string, data dto.ProductsResponseWithPage) error
	GetProductById(uuid uuid.UUID) (dto.ProductResponse, error)
	CreateProductById(key uuid.UUID, data dto.ProductResponse) error
	DeleteAll(key string) error
	Delete(key string) error
}
