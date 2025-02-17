package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"product-api-go/internal/handler/category/dto"
)

type CategoriesRedisMock struct {
	mock.Mock
}

func (m *CategoriesRedisMock) GetAllCategory(key string) (dto.CategoriesResponseWithPage, error) {
	args := m.Called(key)

	var response dto.CategoriesResponseWithPage
	if args.Get(0) != nil {
		response = args.Get(0).(dto.CategoriesResponseWithPage)
	}

	return response, args.Error(1)
}

func (m *CategoriesRedisMock) CreateAllCategory(key string, data dto.CategoriesResponseWithPage) error {
	args := m.Called(key, data)
	return args.Error(0)
}

func (m *CategoriesRedisMock) GetCategoryById(id uuid.UUID) (dto.CategoryResponse, error) {
	args := m.Called(id)
	return args.Get(0).(dto.CategoryResponse), args.Error(1)
}

func (m *CategoriesRedisMock) CreateCategoryById(id uuid.UUID, data dto.CategoryResponse) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *CategoriesRedisMock) DeleteAll(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *CategoriesRedisMock) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}
