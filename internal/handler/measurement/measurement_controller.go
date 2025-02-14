package measurement

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"product-api-go/internal/handler/measurement/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	"product-api-go/internal/pkg/utils"
)

type MeasurementsController struct {
	measurementsService MeasurementsService
}

func NewMeasurementsController(cs MeasurementsService) *MeasurementsController {
	return &MeasurementsController{
		measurementsService: cs,
	}
}

func (mc *MeasurementsController) FindAll(c *gin.Context) {
	var query dto.MeasurementQuery
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

	res, code, err := mc.measurementsService.FindAll(query)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.SuccessPaginate(c, code, "successfully get all measurements", res.Data, res.Page)
}

func (mc *MeasurementsController) FindById(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.MeasurementNotFound, nil)
		return
	}

	res, code, err := mc.measurementsService.FindById(uuid)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, code, "successfully get measurement by id", res)
}

func (mc *MeasurementsController) Create(c *gin.Context) {
	var req dto.MeasurementInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidJsonPayload, err.Error())
		return
	}

	validate := utils.FormatValidationError(req)
	if len(validate) > 0 {
		response.Error(c, http.StatusUnprocessableEntity, constant.MeasurementUnprocessableEntity, validate)
		return
	}

	data, code, err := mc.measurementsService.Create(req)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, 200, "successfully created measurement", data)
}

func (mc *MeasurementsController) Update(c *gin.Context) {
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, constant.MeasurementNotFound, nil)
		return
	}

	var req dto.MeasurementInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, constant.InvalidJsonPayload, err.Error())
		return
	}

	validate := utils.FormatValidationError(req)
	if len(validate) > 0 {
		response.Error(c, http.StatusUnprocessableEntity, constant.MeasurementUnprocessableEntity, validate)
		return
	}

	data, code, err := mc.measurementsService.Update(uuid, req)
	if err != nil {
		response.Error(c, code, err.Error(), nil)
		return
	}

	response.Success(c, 200, "successfully updated measurement", data)
}
