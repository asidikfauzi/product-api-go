package category

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
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
	var query dto.CategoriesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidQueryParameters, err.Error())
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	res, code, err := cc.categoriesService.FindAll(query)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.SuccessPaginate(c, code, "berhasil mendapatkan semua kategori", res.Data, res.Page)
}

func (cc *CategoriesController) FindById(c *gin.Context) {
	id := c.Param("id")

	fmt.Println(id)

	uuid, err := uuid.Parse(id)
	fmt.Println(uuid)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.CategoryNotFound, nil)
		return
	}

	res, code, err := cc.categoriesService.FindById(uuid)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, code, "berhasil mendapatkan kategori", res)
}
