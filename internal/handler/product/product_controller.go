package product

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"product-api-go/internal/handler/product/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	"product-api-go/internal/pkg/utils"
)

type ProductsController struct {
	productsService ProductsService
}

func NewProductsController(cs ProductsService) *ProductsController {
	return &ProductsController{
		productsService: cs,
	}
}

func (cc *ProductsController) FindAll(c *gin.Context) {
	var query dto.ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidQueryParameters, err.Error())
		return
	}

	res, code, err := cc.productsService.FindAll(query)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.SuccessPaginate(c, code, "successfully get all products", res.Data, res.Page)
}

func (cc *ProductsController) FindById(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.ProductNotFound, nil)
		return
	}

	res, code, err := cc.productsService.FindById(uuid)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, code, "successfully get product by id", res)
}

func (cc *ProductsController) Create(c *gin.Context) {
	var req dto.ProductInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidJsonPayload, err.Error())
		return
	}

	validate := utils.FormatValidationError(req)
	if len(validate) > 0 {
		response.Error(c, http.StatusUnprocessableEntity, constant.ProductUnprocessableEntity, validate)
		return
	}

	_, err := uuid.Parse(req.MeasurementID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.MeasurementNotFound, nil)
		return
	}

	for _, category := range req.Categories {
		_, err := uuid.Parse(category)
		if err != nil {
			response.Error(c, http.StatusNotFound, constant.SomeCategoryNotFound, nil)
			return
		}
	}

	data, code, err := cc.productsService.Create(req)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, 200, "successfully created product", data)
}

func (cc *ProductsController) Update(c *gin.Context) {
	id := c.Param("id")

	uuId, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.ProductNotFound, nil)
		return
	}

	var req dto.ProductInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidJsonPayload, err.Error())
		return
	}

	validate := utils.FormatValidationError(req)
	if len(validate) > 0 {
		response.Error(c, http.StatusUnprocessableEntity, constant.ProductUnprocessableEntity, validate)
		return
	}

	_, err = uuid.Parse(req.MeasurementID)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.MeasurementNotFound, nil)
		return
	}

	for _, category := range req.Categories {
		_, err := uuid.Parse(category)
		if err != nil {
			response.Error(c, http.StatusNotFound, constant.SomeCategoryNotFound, nil)
			return
		}
	}

	data, code, err := cc.productsService.Update(uuId, req)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, 200, "successfully updated product", data)
}

func (cc *ProductsController) Delete(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.ProductNotFound, nil)
		return
	}

	data, code, err := cc.productsService.Delete(uuid)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, 200, "successfully deleted product", data)
}
