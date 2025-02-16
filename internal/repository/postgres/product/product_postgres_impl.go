package product

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"product-api-go/internal/handler/product/dto"
	"product-api-go/internal/model"
	"product-api-go/internal/pkg/constant"
	"time"
)

const ErrFormat = "%w: %v"
const FindActiveProductQuery = "id = ? AND deleted_at IS NULL"

type products struct {
	DB *gorm.DB
}

func NewProductsPostgres(db *gorm.DB) ProductsPostgres {
	return &products{
		DB: db,
	}
}

func (c *products) FindAll(q dto.ProductQuery) (res []model.Products, totalItem int64, err error) {
	if q.OrderBy == "" {
		q.OrderBy = "created_at"
	}

	if q.Direction != "ASC" && q.Direction != "DESC" {
		q.Direction = "DESC"
	}

	offset := (q.Page - 1) * q.Limit

	query := c.DB.Model(&model.Products{}).
		Select("products.*, ARRAY_AGG(categories.name) AS category_names").
		Where("products.deleted_at IS NULL").
		Joins("LEFT JOIN product_measurements ON product_measurements.id = products.product_measurement_id").
		Joins("LEFT JOIN category_has_products ON category_has_products.products_id = products.id").
		Joins("LEFT JOIN categories ON categories.id = category_has_products.categories_id").
		Group("products.id, product_measurements.id").
		Preload("ProductMeasurement").
		Preload("Categories")

	if q.Search != "" {
		query = query.Where(
			"products.name ILIKE ? OR product_measurements.name ILIKE ? OR categories.name ILIKE ?",
			"%"+q.Search+"%", "%"+q.Search+"%", "%"+q.Search+"%",
		)
	}

	if q.Category != "" {
		query = query.Where("categories.id = ?", q.Category)
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

	if err = query.Offset(offset).
		Order(fmt.Sprintf("%s %s", q.OrderBy, q.Direction)).
		Find(&res).Error; err != nil {
		return nil, 0, err
	}

	return res, totalItem, nil
}

func (c *products) FindById(id uuid.UUID) (res model.Products, err error) {
	err = c.DB.Where(FindActiveProductQuery, id).
		Preload("ProductMeasurement").
		Preload("Categories").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, constant.ProductNotFound
	}

	return res, err
}

func (c *products) FindByName(name string) (res model.Products, err error) {
	err = c.DB.Where("name = ? AND deleted_at IS NULL", name).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, constant.ProductNotFound
	}

	return res, err
}

func (c *products) FindByNameExcludeID(name string, excludeID uuid.UUID) (res model.Products, err error) {
	err = c.DB.Where("name = ? AND id != ? AND deleted_at IS NULL", name, excludeID).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, constant.ProductNotFound
	}

	return res, err
}

func (c *products) Create(input dto.ProductInput) (res model.Products, err error) {
	tx := c.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf(ErrFormat, constant.UnexpectedError, r)
		}
	}()

	measurementUUID, _ := uuid.Parse(input.MeasurementID)

	product := model.Products{
		Name:                 input.Name,
		Description:          input.Description,
		PurchasePrice:        input.PurchasePrice,
		SellingPrice:         input.SellingPrice,
		TotalStock:           input.TotalStock,
		MinimumStock:         input.MinimumStock,
		Image:                input.Image,
		ProductMeasurementID: measurementUUID,
	}

	if err = tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return res, err
	}

	if len(input.Categories) > 0 {
		var categories []model.Categories
		if err = tx.Where("id IN ?", input.Categories).Find(&categories).Error; err != nil {
			tx.Rollback()
			return res, err
		}

		if len(categories) != len(input.Categories) {
			tx.Rollback()
			return res, constant.SomeCategoryNotFound
		}

		if err = tx.Model(&product).Association("Categories").Replace(categories); err != nil {
			tx.Rollback()
			return res, err
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return res, err
	}

	return product, nil
}

func (c *products) Update(id uuid.UUID, input dto.ProductInput) (res model.Products, err error) {
	tx := c.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf(ErrFormat, constant.UnexpectedError, r)
		}
	}()

	res, err = c.FindById(id)
	if err != nil {
		tx.Rollback()
		return res, err
	}

	measurementUUID, _ := uuid.Parse(input.MeasurementID)

	res.Name = input.Name
	res.Description = input.Description
	res.PurchasePrice = input.PurchasePrice
	res.SellingPrice = input.SellingPrice
	res.TotalStock = input.TotalStock
	res.MinimumStock = input.MinimumStock
	res.Image = input.Image
	res.ProductMeasurementID = measurementUUID

	if err = tx.Save(&res).Error; err != nil {
		tx.Rollback()
		return res, err
	}

	if len(input.Categories) > 0 {
		var categories []model.Categories
		if err = tx.Where("id IN ?", input.Categories).Find(&categories).Error; err != nil {
			tx.Rollback()
			return res, err
		}

		if len(categories) != len(input.Categories) {
			tx.Rollback()
			return res, constant.SomeCategoryNotFound
		}

		if err = tx.Model(&res).Association("Categories").Replace(categories); err != nil {
			tx.Rollback()
			return res, err
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return res, err
	}

	return res, nil
}

func (c *products) Delete(id uuid.UUID) (res model.Products, err error) {
	tx := c.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf(ErrFormat, constant.UnexpectedError, r)
		}
	}()

	res, err = c.FindById(id)
	if err != nil {
		return res, err
	}

	if err := tx.Model(&res).Update("deleted_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return res, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return res, err
	}

	return res, nil
}
