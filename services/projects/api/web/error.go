package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/useurmind/kubelab/services/projects/api/models"
)

func NewApiError(errorType string, msg string, params ...interface{}) models.ApiError {
	return models.ApiError{
		Error: fmt.Sprintf(msg, params...),
	}
}

func AbortWithValidationError(c *gin.Context, msg string, validationData map[string]string) {
	apiError := models.ApiError{
		Type: models.ErrorTypeValidationError,
		Error: msg,
		Data: validationData,
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, apiError)
}

func AbortWithNotFoundError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(http.StatusNotFound, NewApiError(models.ErrorTypeNotFound, "Not found"))
}

func AbortWithInternalError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, NewApiError(models.ErrorTypeInternalError, "Internal server error"))
}

func AbortWithBadRequest(c *gin.Context, message string, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, NewApiError(models.ErrorTypeBadRequest, message))
}