package measurement

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"product-api-go/internal/handler/measurement/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	"product-api-go/internal/repository/postgres/measurement"
)

type measurementsService struct {
	measurementsPostgres measurement.MeasurementsPostgres
}

func NewMeasurementsService(cp measurement.MeasurementsPostgres) MeasurementsService {
	return &measurementsService{
		measurementsPostgres: cp,
	}
}

func (c *measurementsService) FindAll(q dto.MeasurementQuery) (res dto.MeasurementsResponseWithPage, code int, err error) {
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

	return res, http.StatusOK, nil
}

func (c *measurementsService) FindById(id uuid.UUID) (res dto.MeasurementResponse, code int, err error) {
	measurement, err := c.measurementsPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if measurement.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.MeasurementNotFound)
	}

	res.ID = measurement.ID
	res.Name = measurement.Name

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

	checkNameExists, err := c.measurementsPostgres.FindByName(input.Name)
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

	return res, http.StatusOK, nil
}
