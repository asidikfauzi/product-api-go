package category

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/model"
)

type categories struct {
	DB *gorm.DB
}

func NewCategoriesPostgres(db *gorm.DB) CategoriesPostgres {
	return &categories{
		DB: db,
	}
}

func (c *categories) FindAll(q dto.CategoriesQuery) (res []model.Categories, totalItem int64, err error) {
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
	if err := c.DB.Where("id = ?", id).First(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}
