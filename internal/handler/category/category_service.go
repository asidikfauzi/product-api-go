package category

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	"product-api-go/internal/pkg/utils"
	postgres "product-api-go/internal/repository/postgres/category"
	redis "product-api-go/internal/repository/redis/category"
)

type categoriesService struct {
	categoriesPostgres postgres.CategoriesPostgres
	categoriesRedis    redis.CategoriesRedis
}

func NewCategoriesService(cp postgres.CategoriesPostgres, cr redis.CategoriesRedis) CategoriesService {
	return &categoriesService{
		categoriesPostgres: cp,
		categoriesRedis:    cr,
	}
}

func normalizeCategoryQuery(q dto.CategoryQuery) dto.CategoryQuery {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 10
	}
	if q.Paginate == "" || (q.Paginate != "false" && q.Paginate != "true") {
		q.Paginate = "true"
	}
	return q
}

func (c *categoriesService) FindAll(q dto.CategoryQuery) (res dto.CategoriesResponseWithPage, code int, err error) {
	q = normalizeCategoryQuery(q)

	if q.Paginate == "false" && q.Search == "" {
		getCache, err := c.categoriesRedis.GetAllCategory(fmt.Sprintf(constant.AllCategoriesKey, q.OrderBy, q.Direction))
		if err != nil {
			log.Printf(err.Error())
		}

		if getCache.Data != nil {
			return getCache, http.StatusOK, nil
		}
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

	if q.Paginate == "false" && q.Search == "" {
		err := c.categoriesRedis.CreateAllCategory(fmt.Sprintf(constant.AllCategoriesKey, q.OrderBy, q.Direction), res)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	return res, http.StatusOK, nil
}

func (c *categoriesService) FindById(id uuid.UUID) (res dto.CategoryResponse, code int, err error) {
	getCache, err := c.categoriesRedis.GetCategoryById(id)
	if getCache.ID != uuid.Nil {
		return getCache, http.StatusOK, nil
	}

	if !errors.Is(err, constant.KeyRedisNotExists) {
		log.Printf("Redis error: %v", err)
	}

	category, err := c.categoriesPostgres.FindById(id)
	if err != nil {
		if errors.Is(err, constant.CategoryNotFound) {
			return res, http.StatusNotFound, constant.CategoryNotFound
		}
		return res, http.StatusInternalServerError, err
	}

	res.ID = category.ID
	res.Name = category.Name
	res.CreatedAt = utils.FormatTime(category.CreatedAt)
	res.CreatedBy = &category.CreatedBy
	res.UpdatedAt = utils.FormatTime(category.UpdatedAt)
	res.UpdatedBy = &category.UpdatedBy

	c.categoriesRedis.CreateCategoryById(id, res)

	return res, http.StatusOK, nil
}

func (c *categoriesService) Create(input dto.CategoryInput) (res dto.CategoryResponse, code int, err error) {
	_, err = c.categoriesPostgres.FindByName(input.Name)
	if err == nil {
		return res, http.StatusConflict, constant.CategoryAlreadyExists
	}

	newCategory, err := c.categoriesPostgres.Create(input)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.CategoryResponse{
		ID:   newCategory.ID,
		Name: newCategory.Name,
	}

	c.categoriesRedis.DeleteAll(constant.DeleteAllCategoryKey)

	return res, http.StatusCreated, nil
}

func (c *categoriesService) Update(id uuid.UUID, input dto.CategoryInput) (res dto.CategoryResponse, code int, err error) {
	_, err = c.categoriesPostgres.FindById(id)
	if err != nil {
		if errors.Is(err, constant.CategoryNotFound) {
			return res, http.StatusNotFound, constant.CategoryNotFound
		}
		return res, http.StatusInternalServerError, err
	}

	_, err = c.categoriesPostgres.FindByNameExcludeID(input.Name, id)
	if err == nil {
		return res, http.StatusConflict, constant.CategoryAlreadyExists
	}

	editCategory, err := c.categoriesPostgres.Update(id, input)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.CategoryResponse{
		ID:   editCategory.ID,
		Name: editCategory.Name,
	}

	c.categoriesRedis.DeleteAll(constant.DeleteAllCategoryKey)
	c.categoriesRedis.Delete(fmt.Sprintf(constant.CategoryByIdKey, res.ID))

	return res, http.StatusOK, nil
}

func (c *categoriesService) Delete(id uuid.UUID) (res dto.CategoryResponse, code int, err error) {
	_, err = c.categoriesPostgres.FindById(id)
	if err != nil {
		if errors.Is(err, constant.CategoryNotFound) {
			return res, http.StatusNotFound, constant.CategoryNotFound
		}
		return res, http.StatusInternalServerError, err
	}

	deleteCategory, err := c.categoriesPostgres.Delete(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.CategoryResponse{
		ID:   deleteCategory.ID,
		Name: deleteCategory.Name,
	}

	c.categoriesRedis.DeleteAll(constant.DeleteAllCategoryKey)
	c.categoriesRedis.Delete(fmt.Sprintf(constant.CategoryByIdKey, res.ID))

	return res, http.StatusOK, nil
}
