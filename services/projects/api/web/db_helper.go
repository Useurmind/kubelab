package web

import (
	"github.com/gin-gonic/gin"
	"github.com/useurmind/kubelab/services/projects/api/repository")


func UseDBContext(c *gin.Context, dbSystem repository.DBSystem, useContext func(c *gin.Context, dbContext repository.DBContext)) {
	dbContext := dbSystem.NewContext()
	defer dbContext.Close()

	useContext(c, dbContext)
}

// CommitOnSuccess commits the db context if the function returns true.
// It rolls back if false is returned.
func CommitOnSuccess(c *gin.Context, dbContext repository.DBContext, f func() bool) {
	result := f()

	if result {
		err := dbContext.Commit()
		if err != nil {
			c.Error(err)
			return
		}
	} else {
		err := dbContext.Rollback()
		if err != nil {
			c.Error(err)
			return
		}
	}
}