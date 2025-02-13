package response

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func Error(c *gin.Context, code int, message string, details interface{}) {
	resp := ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}

	c.JSON(code, resp)
}
