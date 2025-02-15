package dto

type ProductInput struct {
	Name          string   `json:"name" validate:"required"`
	Description   *string  `json:"description"`
	PurchasePrice int      `json:"purchase_price" validate:"required"`
	SellingPrice  int      `json:"selling_price" validate:"required"`
	TotalStock    int      `json:"total_stock" validate:"required"`
	MinimumStock  int      `json:"minimum_stock" validate:"required"`
	Image         *string  `json:"image" validate:"required"`
	Categories    []string `json:"categories" validate:"required"`
	MeasurementID string   `json:"measurement_id" validate:"required"`
}
