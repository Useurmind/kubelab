package web

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)


func getIntParam(c *gin.Context, name string) (int64, error) {
	if c.Param(name) == "" {
		errMsg := "Parameter " + name + " expected but not set"
		AbortWithBadRequest(c, errMsg, nil)
		return 0, fmt.Errorf(errMsg)
	}

	value, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil {
		AbortWithBadRequest(c, "Could not parse input parameter "+name, err)
		return 0, err
	}

	return value, nil
}

func getIntParamOrDefault(c *gin.Context, name string, defaultValue int64) (int64, error) {
	var err error = nil
	value := defaultValue
	if c.Param(name) != "" {
		value, err = getIntParam(c, name)
		if err != nil {
			return 0, err
		}
	}
	return value, nil
}