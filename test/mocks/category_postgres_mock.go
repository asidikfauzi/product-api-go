package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/model"
	"product-api-go/internal/pkg/constant"
)

type CategoriesPostgresMock struct {
	mock.Mock
}

func (m *CategoriesPostgresMock) FindAll(q dto.CategoryQuery) ([]model.Categories, int64, error) {
	args := m.Called(q)

	var categories []model.Categories
	if args.Get(0) != nil {
		categories = args.Get(0).([]model.Categories)
	}

	var totalItems int64
	if args.Get(1) != nil {
		totalItems = args.Get(1).(int64)
	}

	return categories, totalItems, args.Error(2)
}

func (m *CategoriesPostgresMock) FindById(id uuid.UUID) (model.Categories, error) {
	args := m.Called(id)
	if category, ok := args.Get(0).(model.Categories); ok {
		return category, args.Error(1)
	}
	return model.Categories{}, constant.CategoryNotFound
}

func (m *CategoriesPostgresMock) FindManyById(ids []uuid.UUID) ([]model.Categories, error) {
	args := m.Called(ids)
	return args.Get(0).([]model.Categories), args.Error(1)
}

func (m *CategoriesPostgresMock) FindByName(name string) (model.Categories, error) {
	args := m.Called(name)
	return args.Get(0).(model.Categories), args.Error(1)
}

func (m *CategoriesPostgresMock) FindByNameExcludeID(name string, id uuid.UUID) (model.Categories, error) {
	args := m.Called(name, id)
	return args.Get(0).(model.Categories), args.Error(1)
}

func (m *CategoriesPostgresMock) Create(input dto.CategoryInput) (model.Categories, error) {
	args := m.Called(input)
	return args.Get(0).(model.Categories), args.Error(1)
}

func (m *CategoriesPostgresMock) Update(id uuid.UUID, input dto.CategoryInput) (model.Categories, error) {
	args := m.Called(id, input)
	return args.Get(0).(model.Categories), args.Error(1)
}

func (m *CategoriesPostgresMock) Delete(id uuid.UUID) (model.Categories, error) {
	args := m.Called(id)
	return args.Get(0).(model.Categories), args.Error(1)
}
