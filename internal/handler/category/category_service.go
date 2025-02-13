package category

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	"product-api-go/internal/pkg/utils"
	"product-api-go/internal/repository/postgres/category"
	"time"
)

type categoriesService struct {
	categoriesPostgres category.CategoriesPostgres
}

func NewCategoriesService(cp category.CategoriesPostgres) CategoriesService {
	return &categoriesService{
		categoriesPostgres: cp,
	}
}

func (c *categoriesService) FindAll(q dto.CategoriesQuery) (res dto.FindAllCategoriesResponse, code int, err error) {
	categories, totalItems, err := c.categoriesPostgres.FindAll(q)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	for _, category := range categories {
		res.Data = append(res.Data, dto.FindAllCategory{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	res.Page = response.PaginationResponse{
		TotalItems:   totalItems,
		ItemCount:    len(res.Data),
		ItemsPerPage: q.Limit,
		CurrentPage:  q.Page,
	}

	return res, http.StatusOK, nil
}

func (c *categoriesService) FindById(id uuid.UUID) (res dto.FindByIdCategoryResponse, code int, err error) {
	category, err := c.categoriesPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if category.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.CategoryNotFound)
	}

	formatTime := func(t time.Time) string {
		formattedTime, _ := utils.FormatTimeWithTimezone(t)
		return formattedTime
	}

	res.ID = category.ID
	res.Name = category.Name
	res.CreatedAt = formatTime(category.CreatedAt)
	res.CreatedBy = category.CreatedBy
	res.UpdatedAt = formatTime(category.UpdatedAt)
	res.UpdatedBy = category.UpdatedBy

	return res, http.StatusOK, nil
}
