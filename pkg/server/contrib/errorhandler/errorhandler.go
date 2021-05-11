package errorhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type berrypostErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
}

func JSONErrorHandler() gin.HandlerFunc {
	return jsonErrorHandlerT(gin.ErrorTypeAny)
}

func jsonErrorHandlerT(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)
		if len(detectedErrors) <= 0 {
			return
		}
		response := &berrypostErrorResponse{
			Code:    "internal_server_error",
			Message: "Internal Server Error",
			Detail:  detectedErrors.JSON(),
		}
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
	}
}
