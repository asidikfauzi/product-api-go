package category

import (
	"github.com/google/uuid"
	"product-api-go/internal/handler/category/dto"
)

type CategoriesRedis interface {
	GetAllCategory(key string) (dto.CategoriesResponseWithPage, error)
	CreateAllCategory(key string, data dto.CategoriesResponseWithPage) error
	GetCategoryById(key uuid.UUID) (dto.CategoryResponse, error)
	CreateCategoryById(key uuid.UUID, data dto.CategoryResponse) error
	DeleteAll(key string) error
	Delete(key string) error
}
