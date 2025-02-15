package product

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/product/dto"
	"product-api-go/internal/model"
)

type ProductsPostgres interface {
	FindAll(q dto.ProductQuery) ([]model.Products, int64, error)
	FindById(id uuid.UUID) (model.Products, error)
	FindByName(name string) (model.Products, error)
	FindByNameExcludeID(name string, excludeID uuid.UUID) (res model.Products, err error)
	Create(input dto.ProductInput) (model.Products, error)
	Update(id uuid.UUID, input dto.ProductInput) (model.Products, error)
	Delete(id uuid.UUID) (model.Products, error)
}
