package measurement

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"product-api-go/internal/handler/measurement/dto"
	"product-api-go/internal/model"
	"product-api-go/internal/pkg/constant"
)

type measurements struct {
	DB *gorm.DB
}

func NewMeasurementsPostgres(db *gorm.DB) MeasurementsPostgres {
	return &measurements{
		DB: db,
	}
}

func (c *measurements) FindAll(q dto.MeasurementQuery) (res []model.ProductMeasurements, totalItem int64, err error) {
	if q.OrderBy == "" {
		q.OrderBy = "name"
	}

	if q.Direction != "ASC" && q.Direction != "DESC" {
		q.Direction = "DESC"
	}

	offset := (q.Page - 1) * q.Limit

	query := c.DB.Model(&model.ProductMeasurements{})

	if q.Search != "" {
		query = query.Where("name ILIKE ?", "%"+q.Search+"%")
	}

	if err = query.Count(&totalItem).Error; err != nil {
		return nil, 0, err
	}

	if q.Paginate == "true" {
		if q.Limit > 0 {
			query = query.Limit(q.Limit)
		}

		query = query.Offset(offset)
	}

	if err = query.Order(fmt.Sprintf("%s %s", q.OrderBy, q.Direction)).Find(&res).Error; err != nil {
		return nil, 0, err
	}

	return res, totalItem, nil
}

func (c *measurements) FindById(id uuid.UUID) (res model.ProductMeasurements, err error) {
	err = c.DB.Where("id = ?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, constant.MeasurementNotFound
	}

	return res, err
}

func (c *measurements) FindByName(name string) (res model.ProductMeasurements, err error) {
	err = c.DB.Where("name = ?", name).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, constant.MeasurementNotFound
	}

	return res, err
}

func (c *measurements) FindByNameExcludeID(name string, excludeID uuid.UUID) (res model.ProductMeasurements, err error) {
	err = c.DB.Where("name = ? AND id != ?", name, excludeID).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, constant.MeasurementNotFound
	}

	return res, err
}

func (c *measurements) Create(input dto.MeasurementInput) (res model.ProductMeasurements, err error) {
	measurement := model.ProductMeasurements{
		Name: input.Name,
	}

	if err = c.DB.Model(&res).Create(&measurement).Error; err != nil {
		return res, err
	}

	return measurement, nil
}

func (c *measurements) Update(id uuid.UUID, input dto.MeasurementInput) (res model.ProductMeasurements, err error) {
	res, err = c.FindById(id)
	if err != nil {
		return res, err
	}

	updateData := model.ProductMeasurements{
		Name: input.Name,
	}

	if err = c.DB.Model(&res).Updates(updateData).Error; err != nil {
		return res, err
	}

	return res, nil
}
