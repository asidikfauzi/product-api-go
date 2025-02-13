package category

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/model"
)

type CategoriesPostgres interface {
	FindAll(q dto.CategoriesQuery) ([]model.Categories, int64, error)
	FindById(id uuid.UUID) (model.Categories, error)
}
