package dto

type ProductQuery struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Search    string `form:"search"`
	Category  string `form:"category"`
	OrderBy   string `form:"order_by"`
	Direction string `form:"direction"`
}
