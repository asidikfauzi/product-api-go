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
)

type categoriesService struct {
	categoriesPostgres category.CategoriesPostgres
}

func NewCategoriesService(cp category.CategoriesPostgres) CategoriesService {
	return &categoriesService{
		categoriesPostgres: cp,
	}
}

func (c *categoriesService) FindAll(q dto.CategoryQuery) (res dto.CategoriesResponseWithPage, code int, err error) {
	if q.Page == 0 {
		q.Page = 1
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	if q.Paginate == "" || q.Paginate != "false" && q.Paginate != "true" {
		q.Paginate = "true"
	}

	categories, totalItems, err := c.categoriesPostgres.FindAll(q)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	for _, category := range categories {
		res.Data = append(res.Data, dto.CategoryResponse{
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

func (c *categoriesService) FindById(id uuid.UUID) (res dto.CategoryResponse, code int, err error) {
	category, err := c.categoriesPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if category.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.CategoryNotFound)
	}

	res.ID = category.ID
	res.Name = category.Name
	res.CreatedAt = utils.FormatTime(category.CreatedAt)
	res.CreatedBy = &category.CreatedBy
	res.UpdatedAt = utils.FormatTime(category.UpdatedAt)
	res.UpdatedBy = &category.UpdatedBy

	return res, http.StatusOK, nil
}

func (c *categoriesService) Create(input dto.CategoryInput) (res dto.CategoryResponse, code int, err error) {
	checkIsExists, err := c.categoriesPostgres.FindByName(input.Name)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkIsExists.ID != uuid.Nil {
		return res, http.StatusConflict, errors.New(constant.CategoryAlreadyExists)
	}

	newCategory, err := c.categoriesPostgres.Create(input)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.CategoryResponse{
		ID:   newCategory.ID,
		Name: newCategory.Name,
	}

	return res, http.StatusCreated, nil
}

func (c *categoriesService) Update(id uuid.UUID, input dto.CategoryInput) (res dto.CategoryResponse, code int, err error) {
	checkIsExists, err := c.categoriesPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkIsExists.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.CategoryNotFound)
	}

	checkNameExists, err := c.categoriesPostgres.FindByNameExcludeID(input.Name, id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkNameExists.ID != uuid.Nil {
		return res, http.StatusConflict, errors.New(constant.CategoryAlreadyExists)
	}

	editCategory, err := c.categoriesPostgres.Update(id, input)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.CategoryResponse{
		ID:   editCategory.ID,
		Name: editCategory.Name,
	}

	return res, http.StatusOK, nil
}

func (c *categoriesService) Delete(id uuid.UUID) (res dto.CategoryResponse, code int, err error) {
	checkIsExists, err := c.categoriesPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkIsExists.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.CategoryNotFound)
	}

	deleteCategory, err := c.categoriesPostgres.Delete(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.CategoryResponse{
		ID:   deleteCategory.ID,
		Name: deleteCategory.Name,
	}

	return res, http.StatusOK, nil
}
