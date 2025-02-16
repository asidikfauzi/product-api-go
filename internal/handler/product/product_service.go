package product

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	dtoCat "product-api-go/internal/handler/category/dto"
	"product-api-go/internal/handler/product/dto"
	"product-api-go/internal/pkg/constant"
	"product-api-go/internal/pkg/response"
	"product-api-go/internal/pkg/utils"
	"product-api-go/internal/repository/postgres/category"
	"product-api-go/internal/repository/postgres/measurement"
	postgres "product-api-go/internal/repository/postgres/product"
	redis "product-api-go/internal/repository/redis/product"
)

type productsService struct {
	productsPostgres    postgres.ProductsPostgres
	productRedis        redis.ProductsRedis
	measurementPostgres measurement.MeasurementsPostgres
	categoryPostgres    category.CategoriesPostgres
}

func NewProductsService(cp postgres.ProductsPostgres, pr redis.ProductsRedis, mp measurement.MeasurementsPostgres, cap category.CategoriesPostgres) ProductsService {
	return &productsService{
		productsPostgres:    cp,
		productRedis:        pr,
		measurementPostgres: mp,
		categoryPostgres:    cap,
	}
}

func normalizeProductQuery(q dto.ProductQuery) dto.ProductQuery {
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

func (c *productsService) FindAll(q dto.ProductQuery) (res dto.ProductsResponseWithPage, code int, err error) {
	q = normalizeProductQuery(q)

	if q.Paginate == "false" && q.Search == "" {
		getCache, err := c.productRedis.GetAllProduct(fmt.Sprintf(constant.AllProductsKey, q.OrderBy, q.Direction, q.Category))
		if err != nil {
			log.Printf(err.Error())
		}

		if getCache.Data != nil {
			return getCache, http.StatusOK, nil
		}
	}

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

	if q.Paginate == "false" && q.Search == "" {
		err := c.productRedis.CreateAllProduct(fmt.Sprintf(constant.AllProductsKey, q.OrderBy, q.Direction, q.Category), res)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	return res, http.StatusOK, nil
}

func (c *productsService) FindById(id uuid.UUID) (res dto.ProductResponse, code int, err error) {
	getCache, err := c.productRedis.GetProductById(id)
	if getCache.ID != uuid.Nil {
		return getCache, http.StatusOK, nil
	}

	if !errors.Is(err, constant.KeyRedisNotExists) {
		log.Printf("Redis error: %v", err)
	}

	product, err := c.productsPostgres.FindById(id)
	if err != nil {
		if errors.Is(err, constant.ProductNotFound) {
			return res, http.StatusNotFound, constant.ProductNotFound
		}
		return res, http.StatusInternalServerError, err
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

	c.productRedis.CreateProductById(id, res)

	return res, http.StatusOK, nil
}

func (c *productsService) Create(input dto.ProductInput) (res dto.ProductResponse, code int, err error) {
	uuidMea, _ := uuid.Parse(input.MeasurementID)

	_, err = c.measurementPostgres.FindById(uuidMea)
	if err != nil {
		if errors.Is(err, constant.MeasurementNotFound) {
			return res, http.StatusNotFound, constant.MeasurementNotFound
		}
		return res, http.StatusInternalServerError, err
	}

	var categoryIDs []uuid.UUID
	for _, id := range input.Categories {
		parsedID, _ := uuid.Parse(id)
		categoryIDs = append(categoryIDs, parsedID)
	}

	_, err = c.categoryPostgres.FindManyById(categoryIDs)
	if err != nil {
		if errors.Is(err, constant.SomeCategoryNotFound) {
			return res, http.StatusNotFound, constant.SomeCategoryNotFound
		}
		return res, http.StatusInternalServerError, err
	}

	_, err = c.productsPostgres.FindByName(input.Name)
	if err == nil {
		return res, http.StatusConflict, constant.ProductAlreadyExists
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

	c.productRedis.DeleteAll(constant.DeleteAllProductKey)

	return res, http.StatusCreated, nil
}

func (c *productsService) Update(id uuid.UUID, input dto.ProductInput) (res dto.ProductResponse, code int, err error) {
	_, err = c.productsPostgres.FindById(id)
	if err != nil {
		if errors.Is(err, constant.ProductNotFound) {
			return res, http.StatusNotFound, constant.ProductNotFound
		}
		return res, http.StatusInternalServerError, err
	}

	uuidMea, _ := uuid.Parse(input.MeasurementID)

	_, err = c.measurementPostgres.FindById(uuidMea)
	if err != nil {
		if errors.Is(err, constant.MeasurementNotFound) {
			return res, http.StatusNotFound, constant.MeasurementNotFound
		}
		return res, http.StatusInternalServerError, err
	}

	var categoryIDs []uuid.UUID
	for _, id := range input.Categories {
		parsedID, _ := uuid.Parse(id)
		categoryIDs = append(categoryIDs, parsedID)
	}

	_, err = c.categoryPostgres.FindManyById(categoryIDs)
	if err != nil {
		if errors.Is(err, constant.SomeCategoryNotFound) {
			return res, http.StatusNotFound, constant.SomeCategoryNotFound
		}
		return res, http.StatusInternalServerError, err
	}

	_, err = c.productsPostgres.FindByNameExcludeID(input.Name, id)
	if err == nil {
		return res, http.StatusConflict, constant.ProductAlreadyExists
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

	c.productRedis.DeleteAll(constant.DeleteAllProductKey)
	c.productRedis.Delete(fmt.Sprintf(constant.ProductByIdKey, res.ID))

	return res, http.StatusOK, nil
}

func (c *productsService) Delete(id uuid.UUID) (res dto.ProductResponse, code int, err error) {
	_, err = c.productsPostgres.FindById(id)
	if err != nil {
		if errors.Is(err, constant.ProductNotFound) {
			return res, http.StatusNotFound, constant.ProductNotFound
		}
		return res, http.StatusInternalServerError, err
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

	c.productRedis.DeleteAll(constant.DeleteAllProductKey)
	c.productRedis.Delete(fmt.Sprintf(constant.ProductByIdKey, res.ID))

	return res, http.StatusOK, nil
}
