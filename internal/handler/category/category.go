package category

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/category/dto"
)

type CategoriesService interface {
	FindAll(query dto.CategoriesQuery) (dto.FindAllCategoriesResponse, int, error)
	FindById(id uuid.UUID) (dto.FindByIdCategoryResponse, int, error)
}
