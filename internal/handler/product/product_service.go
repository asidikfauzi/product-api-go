package product

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	dtoCat "product-api-go/internal/handler/category/dto"
	"product-api-go/internal/handler/product/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	"product-api-go/internal/pkg/utils"
	"product-api-go/internal/repository/postgres/category"
	"product-api-go/internal/repository/postgres/measurement"
	"product-api-go/internal/repository/postgres/product"
)

type productsService struct {
	productsPostgres    product.ProductsPostgres
	measurementPostgres measurement.MeasurementsPostgres
	categoryPostgres    category.CategoriesPostgres
}

func NewProductsService(cp product.ProductsPostgres, mp measurement.MeasurementsPostgres, cap category.CategoriesPostgres) ProductsService {
	return &productsService{
		productsPostgres:    cp,
		measurementPostgres: mp,
		categoryPostgres:    cap,
	}
}

func (c *productsService) FindAll(q dto.ProductQuery) (res dto.ProductsResponseWithPage, code int, err error) {
	products, totalItems, err := c.productsPostgres.FindAll(q)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	for _, product := range products {

		var categories []dtoCat.CategoryResponse
		for _, category := range product.Categories {
			categories = append(categories, dtoCat.CategoryResponse{
				ID:   category.ID,
				Name: category.Name,
			})
		}

		res.Data = append(res.Data, dto.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			PurchasePrice: product.PurchasePrice,
			SellingPrice:  product.SellingPrice,
			TotalStock:    product.TotalStock,
			MinimumStock:  product.MinimumStock,
			Image:         product.Image,
			Measurement:   product.ProductMeasurement.Name,
			Categories:    categories,
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

func (c *productsService) FindById(id uuid.UUID) (res dto.ProductResponse, code int, err error) {
	product, err := c.productsPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if product.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.ProductNotFound)
	}

	var categories []dtoCat.CategoryResponse
	for _, category := range product.Categories {
		categories = append(categories, dtoCat.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	res.ID = product.ID
	res.Name = product.Name
	res.Description = utils.FormatDefaultString(product.Description, "")
	res.PurchasePrice = product.PurchasePrice
	res.SellingPrice = product.SellingPrice
	res.TotalStock = product.TotalStock
	res.MinimumStock = product.MinimumStock
	res.Image = product.Image
	res.Measurement = product.ProductMeasurement.Name
	res.Categories = categories
	res.CreatedAt = utils.FormatTime(product.CreatedAt)
	res.CreatedBy = &product.CreatedBy
	res.UpdatedAt = utils.FormatTime(product.UpdatedAt)
	res.UpdatedBy = &product.UpdatedBy

	return res, http.StatusOK, nil
}

func (c *productsService) Create(input dto.ProductInput) (res dto.ProductResponse, code int, err error) {
	uuidMea, _ := uuid.Parse(input.MeasurementID)

	checkMeaExists, err := c.measurementPostgres.FindById(uuidMea)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkMeaExists.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.MeasurementNotFound)
	}

	var categoryIDs []uuid.UUID
	for _, id := range input.Categories {
		parsedID, _ := uuid.Parse(id)
		categoryIDs = append(categoryIDs, parsedID)
	}

	checkCatSomeExists, err := c.categoryPostgres.FindManyById(categoryIDs)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if len(checkCatSomeExists) == 0 {
		return res, http.StatusNotFound, errors.New(constant.SomeCategoryNotFound)
	}

	checkIsExists, err := c.productsPostgres.FindByName(input.Name)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkIsExists.ID != uuid.Nil {
		return res, http.StatusConflict, errors.New(constant.ProductAlreadyExists)
	}

	newProduct, err := c.productsPostgres.Create(input)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.ProductResponse{
		ID:            newProduct.ID,
		Name:          newProduct.Name,
		Description:   utils.FormatDefaultString(newProduct.Description, ""),
		PurchasePrice: newProduct.PurchasePrice,
		SellingPrice:  newProduct.SellingPrice,
		TotalStock:    newProduct.TotalStock,
		MinimumStock:  newProduct.MinimumStock,
		Image:         newProduct.Image,
	}

	return res, http.StatusCreated, nil
}

func (c *productsService) Update(id uuid.UUID, input dto.ProductInput) (res dto.ProductResponse, code int, err error) {
	checkIsExists, err := c.productsPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkIsExists.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.ProductNotFound)
	}

	uuidMea, _ := uuid.Parse(input.MeasurementID)

	checkMeaExists, err := c.measurementPostgres.FindById(uuidMea)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkMeaExists.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.MeasurementNotFound)
	}

	var categoryIDs []uuid.UUID
	for _, id := range input.Categories {
		parsedID, _ := uuid.Parse(id)
		categoryIDs = append(categoryIDs, parsedID)
	}

	checkCatSomeExists, err := c.categoryPostgres.FindManyById(categoryIDs)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if len(checkCatSomeExists) == 0 {
		return res, http.StatusNotFound, errors.New(constant.SomeCategoryNotFound)
	}

	checkNameExists, err := c.productsPostgres.FindByNameExcludeID(input.Name, id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkNameExists.ID != uuid.Nil {
		return res, http.StatusConflict, errors.New(constant.ProductAlreadyExists)
	}

	editProduct, err := c.productsPostgres.Update(id, input)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.ProductResponse{
		ID:            editProduct.ID,
		Name:          editProduct.Name,
		Description:   utils.FormatDefaultString(editProduct.Description, ""),
		PurchasePrice: editProduct.PurchasePrice,
		SellingPrice:  editProduct.SellingPrice,
		TotalStock:    editProduct.TotalStock,
		MinimumStock:  editProduct.MinimumStock,
		Image:         editProduct.Image,
	}

	return res, http.StatusOK, nil
}

func (c *productsService) Delete(id uuid.UUID) (res dto.ProductResponse, code int, err error) {
	checkIsExists, err := c.productsPostgres.FindById(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	if checkIsExists.ID == uuid.Nil {
		return res, http.StatusNotFound, errors.New(constant.ProductNotFound)
	}

	deleteProduct, err := c.productsPostgres.Delete(id)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = dto.ProductResponse{
		ID:            deleteProduct.ID,
		Name:          deleteProduct.Name,
		Description:   utils.FormatDefaultString(deleteProduct.Description, ""),
		PurchasePrice: deleteProduct.PurchasePrice,
		SellingPrice:  deleteProduct.SellingPrice,
		TotalStock:    deleteProduct.TotalStock,
		MinimumStock:  deleteProduct.MinimumStock,
		Image:         deleteProduct.Image,
	}

	return res, http.StatusOK, nil
}
