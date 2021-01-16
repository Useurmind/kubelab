package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


type ApiError struct {
	Error string `json:"error"`
}

func NewApiError(msg string, params ...interface{}) ApiError {
	return ApiError{
		Error: fmt.Sprintf(msg, params...),
	}
}

func AbortWithNotFoundError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(http.StatusNotFound, NewApiError("Not found"))
}

func AbortWithInternalError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, NewApiError("Internal server error"))
}

func AbortWithBadRequest(c *gin.Context, message string, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, NewApiError(message))
}