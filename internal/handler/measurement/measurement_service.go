package measurement

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"product-api-go/internal/handler/measurement/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	postgres "product-api-go/internal/repository/postgres/measurement"
	redis "product-api-go/internal/repository/redis/measurement"
)

type measurementsService struct {
	measurementsPostgres postgres.MeasurementsPostgres
	measurementsRedis    redis.MeasurementsRedis
}

func NewMeasurementsService(cp postgres.MeasurementsPostgres, mr redis.MeasurementsRedis) MeasurementsService {
	return &measurementsService{
		measurementsPostgres: cp,
		measurementsRedis:    mr,
	}
}

func normalizeMeasurementQuery(q dto.MeasurementQuery) dto.MeasurementQuery {
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

func (c *measurementsService) FindAll(q dto.MeasurementQuery) (res dto.MeasurementsResponseWithPage, code int, err error) {
	q = normalizeMeasurementQuery(q)

	if q.Paginate == "false" && q.Search == "" {
		getCache, err := c.measurementsRedis.GetAllMeasurement(fmt.Sprintf(constant.AllMeasurementsKey, q.OrderBy, q.Direction))
		if err != nil {
			log.Printf(err.Error())
		}

		if getCache.Data != nil {
			return getCache, http.StatusOK, nil
		}
	}

	measurements, totalItems, err := c.measurementsPostgres.FindAll(q)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	for _, measurement := range measurements {
		res.Data = append(res.Data, dto.MeasurementResponse{
			ID:   measurement.ID,
			Name: measurement.Name,
		})
	}

	res.Page = response.PaginationResponse{
		TotalItems:   totalItems,
		ItemCount:    len(res.Data),
		ItemsPerPage: q.Limit,
		CurrentPage:  q.Page,
	}

	if q.Paginate == "false" && q.Search == "" {
		err := c.measurementsRedis.CreateAllMeasurement(fmt.Sprintf(constant.AllMeasurementsKey, q.OrderBy, q.Direction), res)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	return res, http.StatusOK, nil
}

func (c *measurementsService) FindById(id uuid.UUID) (res dto.MeasurementResponse, code int, err error) {
	getCache, err := c.measurementsRedis.GetMeasurementById(id)
	if err != nil {
		log.Printf(err.Error())
	}

	if getCache.ID != uuid.Nil {
		return getCache, http.StatusOK, nil
	}

	measurement, err := c.measurementsPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if measurement.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.MeasurementNotFound)
	}

	res.ID = measurement.ID
	res.Name = measurement.Name

	c.measurementsRedis.CreateMeasurementById(id, res)

	return res, http.StatusOK, nil
}

func (c *measurementsService) Create(input dto.MeasurementInput) (res dto.MeasurementResponse, code int, err error) {
	checkIsExists, err := c.measurementsPostgres.FindByName(input.Name)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkIsExists.ID != uuid.Nil {
		return res, http.StatusConflict, errors.New(constant.MeasurementAlreadyExists)
	}

	newMeasurement, err := c.measurementsPostgres.Create(input)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.MeasurementResponse{
		ID:   newMeasurement.ID,
		Name: newMeasurement.Name,
	}

	c.measurementsRedis.DeleteAll(constant.DeleteAllMeasurementKey)

	return res, http.StatusCreated, nil
}

func (c *measurementsService) Update(id uuid.UUID, input dto.MeasurementInput) (res dto.MeasurementResponse, code int, err error) {
	checkIsExists, err := c.measurementsPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkIsExists.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.MeasurementNotFound)
	}

	checkNameExists, err := c.measurementsPostgres.FindByNameExcludeID(input.Name, id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkNameExists.ID != uuid.Nil {
		return res, http.StatusConflict, errors.New(constant.MeasurementAlreadyExists)
	}

	editMeasurement, err := c.measurementsPostgres.Update(id, input)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.MeasurementResponse{
		ID:   editMeasurement.ID,
		Name: editMeasurement.Name,
	}

	c.measurementsRedis.DeleteAll(constant.DeleteAllMeasurementKey)
	c.measurementsRedis.Delete(fmt.Sprintf(constant.MeasurementByIdKey, res.ID))

	return res, http.StatusOK, nil
}
