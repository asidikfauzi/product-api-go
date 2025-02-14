package category

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/category/dto"
)

type CategoriesService interface {
	FindAll(query dto.CategoryQuery) (dto.CategoriesResponseWithPage, int, error)
	FindById(id uuid.UUID) (dto.CategoryResponse, int, error)
	Create(input dto.CategoryInput) (dto.CategoryResponse, int, error)
	Update(id uuid.UUID, input dto.CategoryInput) (dto.CategoryResponse, int, error)
	Delete(id uuid.UUID) (dto.CategoryResponse, int, error)
}
