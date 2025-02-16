package category

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	"product-api-go/internal/pkg/utils"
)

type CategoriesController struct {
	categoriesService CategoriesService
}

func NewCategoriesController(cs CategoriesService) *CategoriesController {
	return &CategoriesController{
		categoriesService: cs,
	}
}

func (cc *CategoriesController) FindAll(c *gin.Context) {
	var query dto.CategoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidQueryParameters.Error(), err.Error())
		return
	}

	res, code, err := cc.categoriesService.FindAll(query)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.SuccessPaginate(c, code, "successfully get all categories", res.Data, res.Page)
}

func (cc *CategoriesController) FindById(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.CategoryNotFound.Error(), nil)
		return
	}

	res, code, err := cc.categoriesService.FindById(uuid)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, code, "successfully get category by id", res)
}

func (cc *CategoriesController) Create(c *gin.Context) {
	var req dto.CategoryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidJsonPayload.Error(), err.Error())
		return
	}

	validate := utils.FormatValidationError(req)
	if len(validate) > 0 {
		response.Error(c, http.StatusUnprocessableEntity, constant.CategoryUnprocessableEntity.Error(), validate)
		return
	}

	data, code, err := cc.categoriesService.Create(req)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, 200, "successfully created category", data)
}

func (cc *CategoriesController) Update(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.CategoryNotFound.Error(), nil)
		return
	}

	var req dto.CategoryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidJsonPayload.Error(), err.Error())
		return
	}

	validate := utils.FormatValidationError(req)
	if len(validate) > 0 {
		response.Error(c, http.StatusUnprocessableEntity, constant.CategoryUnprocessableEntity.Error(), validate)
		return
	}

	data, code, err := cc.categoriesService.Update(uuid, req)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, 200, "successfully updated category", data)
}

func (cc *CategoriesController) Delete(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.CategoryNotFound.Error(), nil)
		return
	}

	data, code, err := cc.categoriesService.Delete(uuid)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, 200, "successfully deleted category", data)
}
