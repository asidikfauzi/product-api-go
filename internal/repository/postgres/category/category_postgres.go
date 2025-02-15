package category

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/model"
)

type CategoriesPostgres interface {
	FindAll(q dto.CategoryQuery) ([]model.Categories, int64, error)
	FindById(id uuid.UUID) (model.Categories, error)
	FindManyById(ids []uuid.UUID) ([]model.Categories, error)
	FindByName(name string) (model.Categories, error)
	FindByNameExcludeID(name string, excludeID uuid.UUID) (model.Categories, error)
	Create(input dto.CategoryInput) (model.Categories, error)
	Update(id uuid.UUID, input dto.CategoryInput) (model.Categories, error)
	Delete(id uuid.UUID) (model.Categories, error)
}
