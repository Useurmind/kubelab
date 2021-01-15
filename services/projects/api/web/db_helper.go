package web

import (
	"github.com/gin-gonic/gin"
	"github.com/useurmind/kubelab/services/projects/api/repository")


func UseDBContext(c *gin.Context, dbSystem repository.DBSystem, useContext func(c *gin.Context, dbContext repository.DBContext)) {
	dbContext := dbSystem.NewContext()
	defer dbContext.Close()

	useContext(c, dbContext)
}