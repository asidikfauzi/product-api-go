package response

import (
	"github.com/gin-gonic/gin"
	"math"
)

type SuccessResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    interface{}         `json:"data,omitempty"`
	Page    *PaginationResponse `json:"page,omitempty"`
}

type PaginationResponse struct {
	TotalItems   int64 `json:"total_items"`
	ItemCount    int   `json:"item_count"`
	ItemsPerPage int   `json:"items_per_page"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	HasNextPage  bool  `json:"has_next_page"`
	HasPrevPage  bool  `json:"has_prev_page"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	resp := SuccessResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}

	c.JSON(code, resp)
}

func SuccessPaginate(c *gin.Context, code int, message string, data interface{}, page PaginationResponse) {
	totalPages := int(math.Ceil(float64(page.TotalItems) / float64(page.ItemsPerPage)))
	hasNextPage := int(page.CurrentPage) < totalPages
	hasPrevPage := int(page.CurrentPage) > 1

	page.HasNextPage = hasNextPage
	page.HasPrevPage = hasPrevPage
	page.TotalPages = int(totalPages)

	resp := SuccessResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Page:    &page,
	}

	c.JSON(code, resp)
}
