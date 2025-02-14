package category

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/model"
	"time"
)

const FindActiveCategoryQuery = "id = ? AND deleted_at IS NULL"

type categories struct {
	DB *gorm.DB
}

func NewCategoriesPostgres(db *gorm.DB) CategoriesPostgres {
	return &categories{
		DB: db,
	}
}

func (c *categories) FindAll(q dto.CategoryQuery) (res []model.Categories, totalItem int64, err error) {
	if q.OrderBy == "" {
		q.OrderBy = "created_at"
	}

	if q.Direction != "ASC" && q.Direction != "DESC" {
		q.Direction = "DESC"
	}

	offset := (q.Page - 1) * q.Limit

	query := c.DB.Model(&model.Categories{}).Where("deleted_at IS NULL")

	if q.Search != "" {
		query = query.Where("name ILIKE ?", "%"+q.Search+"%")
	}

	if err = query.Count(&totalItem).Error; err != nil {
		return nil, 0, err
	}

	if q.Limit > 0 {
		query = query.Limit(q.Limit)
	}

	if err = query.Offset(offset).
		Order(fmt.Sprintf("%s %s", q.OrderBy, q.Direction)).
		Find(&res).Error; err != nil {
		return nil, 0, err
	}

	return res, totalItem, nil
}

func (c *categories) FindById(id uuid.UUID) (res model.Categories, err error) {
	err = c.DB.Where(FindActiveCategoryQuery, id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}

	return res, nil
}

func (c *categories) FindByName(name string) (res model.Categories, err error) {
	err = c.DB.Where("name = ? AND deleted_at IS NULL", name).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}

	return res, err
}

func (c *categories) Create(input dto.CategoryInput) (res model.Categories, err error) {
	category := model.Categories{
		Name:      input.Name,
		CreatedAt: time.Now(),
	}

	if err = c.DB.Model(&res).Create(&category).Error; err != nil {
		return res, err
	}

	return category, nil
}

func (c *categories) Update(id uuid.UUID, input dto.CategoryInput) (res model.Categories, err error) {
	if err = c.DB.First(&res, FindActiveCategoryQuery, id).Error; err != nil {
		return res, err
	}

	updateData := model.Categories{
		Name:      input.Name,
		UpdatedAt: time.Now(),
	}

	if err = c.DB.Model(&res).Updates(updateData).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (c *categories) Delete(id uuid.UUID) (res model.Categories, err error) {
	if err := c.DB.First(&res, FindActiveCategoryQuery, id).Error; err != nil {
		return res, err
	}

	timeNow := time.Now()
	uuidUser := uuid.Nil

	if err := c.DB.Model(&res).Updates(model.Categories{
		DeletedAt: &timeNow,
		DeletedBy: &uuidUser,
	}).Error; err != nil {
		return res, err
	}

	return res, nil
}
