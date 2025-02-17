package service_test

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"product-api-go/internal/handler/category"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/model"
	"product-api-go/internal/pkg/constant"
	"product-api-go/test/mocks"
	"testing"
)

type CategoriesServiceTestSuite struct {
	suite.Suite
	categoriesPostgres *mocks.CategoriesPostgresMock
	categoriesRedis    *mocks.CategoriesRedisMock
	service            category.CategoriesService
}

func (suite *CategoriesServiceTestSuite) SetupTest() {
	suite.categoriesPostgres = new(mocks.CategoriesPostgresMock)
	suite.categoriesRedis = new(mocks.CategoriesRedisMock)
	suite.service = category.NewCategoriesService(suite.categoriesPostgres, suite.categoriesRedis)
}

func (suite *CategoriesServiceTestSuite) TestFindAll() {
	query := dto.CategoryQuery{
		Paginate:  "false",
		Search:    "",
		Page:      1,
		Limit:     10,
		OrderBy:   "name",
		Direction: "asc",
	}
	key := fmt.Sprintf(constant.AllCategoriesKey, query.OrderBy, query.Direction)

	categoryDataPostgres := []model.Categories{
		{ID: uuid.New(), Name: "Beverages"},
		{ID: uuid.New(), Name: "Snacks"},
	}

	categoryDataRedis := []dto.CategoryResponse{
		{ID: categoryDataPostgres[0].ID, Name: categoryDataPostgres[0].Name},
		{ID: categoryDataPostgres[1].ID, Name: categoryDataPostgres[1].Name},
	}

	// Test when data is cached
	suite.categoriesRedis.
		On("GetAllCategory", key).
		Return(dto.CategoriesResponseWithPage{Data: categoryDataRedis}, nil).
		Once()

	res, code, err := suite.service.FindAll(query)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, code)
	assert.Equal(suite.T(), categoryDataRedis, res.Data)

	// Test when no cache, fetch from DB
	suite.categoriesRedis.
		On("GetAllCategory", key).
		Return(dto.CategoriesResponseWithPage{}, errors.New("cache miss")).
		Once()

	suite.categoriesPostgres.
		On("FindAll", mock.Anything).
		Return(categoryDataPostgres, int64(len(categoryDataPostgres)), nil).
		Once()

	suite.categoriesRedis.
		On("CreateAllCategory", key, mock.Anything).
		Return(nil).
		Once()

	res, code, err = suite.service.FindAll(query)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, code)

	// Konversi categoryDataPostgres ke []dto.CategoryResponse untuk validasi assertion
	expectedResponse := []dto.CategoryResponse{}
	for _, c := range categoryDataPostgres {
		expectedResponse = append(expectedResponse, dto.CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	assert.Equal(suite.T(), expectedResponse, res.Data)
}

func (suite *CategoriesServiceTestSuite) TestFindById() {
	idStr := "d12b9e7b-0a86-475d-b78c-d17ef7eceff5"
	id, _ := uuid.Parse(idStr)
	categoryPostgres := model.Categories{ID: id, Name: "Beverages"}
	categoryRedis := dto.CategoryResponse{ID: id, Name: "Beverages"}

	// Case 1: Data ditemukan di Redis
	suite.categoriesRedis.On("GetCategoryById", id).Return(categoryRedis, nil).Once()
	res, code, err := suite.service.FindById(id)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, code)
	assert.Equal(suite.T(), categoryRedis, res)

	// Case 2: Tidak ada di Redis, lalu ambil dari DB
	suite.categoriesRedis.On("GetCategoryById", id).Return(dto.CategoryResponse{}, constant.KeyRedisNotExists).Once()
	suite.categoriesPostgres.On("FindById", id).Return(categoryPostgres, nil).Once()
	suite.categoriesRedis.On("CreateCategoryById", id, mock.AnythingOfType("dto.CategoryResponse")).Return(nil).Once()

	res, code, err = suite.service.FindById(id)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, code)
	assert.Equal(suite.T(), categoryPostgres.ID, res.ID)
	assert.Equal(suite.T(), categoryPostgres.Name, res.Name)

	// Case 3: Tidak ada di Redis maupun DB
	suite.categoriesRedis.On("GetCategoryById", id).Return(dto.CategoryResponse{}, constant.KeyRedisNotExists).Once()
	suite.categoriesPostgres.On("FindById", id).Return(model.Categories{}, constant.CategoryNotFound).Once()

	res, code, err = suite.service.FindById(id)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNotFound, code)
}

// Run test suite
func TestCategoriesServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CategoriesServiceTestSuite))
}
