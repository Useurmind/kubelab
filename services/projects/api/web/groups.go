package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/useurmind/kubelab/services/projects/api/models"
	"github.com/useurmind/kubelab/services/projects/api/repository"
)

func HandleGroups(basePath string, router *gin.Engine, dbSystem repository.DBSystem) {
	// Create a group
	router.POST(basePath, func(c *gin.Context) { UseDBContext(c, dbSystem, createGroup) })

	// Update a group
	router.PUT(basePath+"/:id", func(c *gin.Context) { UseDBContext(c, dbSystem, updateGroup) })

	// GET group by id
	router.GET(basePath+"/:id", func(c *gin.Context) { UseDBContext(c, dbSystem, getGroupByID) })

	// List groups
	router.GET(basePath, func(c *gin.Context) { UseDBContext(c, dbSystem, listGroups) })

	// DELETE group
	router.DELETE(basePath+"/:id", func(c *gin.Context) { UseDBContext(c, dbSystem, deleteGroup) })
}

func createGroup(c *gin.Context, dbContext repository.DBContext) {
	repo, err := dbContext.GetGroupRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		AbortWithInternalError(c, err)
		return
	}
	group := models.Group{}
	err = json.Unmarshal(jsonBody, &group)
	if err != nil {
		log.Error().Str("jsonbody", string(jsonBody)).Msg("Could not parse json body")
		AbortWithBadRequest(c, "Could not parse body as group object.", err)
		return
	}

	if !group.IsNew() {
		AbortWithBadRequest(c, "Group already has an id, did you intend to update an existing group? Use PUT /groups/:id instead.", nil)
		return
	}

	if group.HasAnyProjects() {
		AbortWithBadRequest(c, "Cannot create a new group with projects in it.", nil)
		return
	}

	CommitOnSuccess(c, dbContext, func() bool {
		groupResult, err := repo.CreateOrUpdate(c, &group)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		c.JSON(200, groupResult)
		return true
	})
}

func updateGroup(c *gin.Context, dbContext repository.DBContext) {
	groupRepo, err := dbContext.GetGroupRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}
	projectRepo, err := dbContext.GetProjectRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	id, err := getIntParam(c, "id")
	if err != nil {
		return
	}

	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		AbortWithInternalError(c, err)
		return
	}

	group := models.Group{}
	err = json.Unmarshal(jsonBody, &group)
	if err != nil {
		log.Error().Str("jsonbody", string(jsonBody)).Msg("Could not parse json body")
		AbortWithBadRequest(c, "Could not parse body as group object.", err)
		return
	}

	if group.IsNew() {
		AbortWithBadRequest(c, fmt.Sprintf("Group id (%d) is zero or not set (path id %d). This indicates a new group. Use POST /groups instead.", group.Id, id), nil)
		return
	}

	if group.Id != id {
		AbortWithBadRequest(c, fmt.Sprintf("Group id %d does not match id %d in url, did you really want to put this group?", group.Id, id), nil)
		return
	}

	CommitOnSuccess(c, dbContext, func() bool {
		projectCount, err := projectRepo.CountByGroupID(c, group.Id)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}
	
		if group.GetProjectCount() != int(projectCount) {
			projects, err := projectRepo.GetByGroupID(c, group.Id)
			if err != nil {
				AbortWithInternalError(c, err)
				return false
			}

			err = group.RefreshProjectRefs(projects)
			
			if err != nil {
				AbortWithInternalError(c, err)
				return false
			}
		}

		groupResult, err := groupRepo.CreateOrUpdate(c, &group)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		c.JSON(200, groupResult)
		return true
	})
}

func getGroupByID(c *gin.Context, dbContext repository.DBContext) {
	repo, err := dbContext.GetGroupRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	id, err := getIntParam(c, "id")
	if err != nil {
		return
	}

	group, err := repo.Get(c, id)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	if group == nil {
		AbortWithNotFoundError(c, nil)
		return
	}

	c.JSON(200, group)
}

func listGroups(c *gin.Context, dbContext repository.DBContext) {
	repo, err := dbContext.GetGroupRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}
	
	startIndex, err := getIntParamOrDefault(c, "start", 0)
	if err != nil {
		return
	}
	count, err := getIntParamOrDefault(c, "count", 50)
	if err != nil {
		return
	}

	groups, err := repo.List(c, startIndex, count)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	c.JSON(200, groups)
}

func deleteGroup(c *gin.Context, dbContext repository.DBContext) {
	groupRepo, err := dbContext.GetGroupRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}
	projectRepo, err := dbContext.GetProjectRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	id, err := getIntParam(c, "id")
	if err != nil {
		return
	}

	CommitOnSuccess(c, dbContext, func() bool {
		projectCount, err := projectRepo.CountByGroupID(c, id)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		if projectCount > 0 {
			AbortWithBadRequest(c, "Cannot deleted group that has still projects assigned to it", nil)
			return false
		}

		err = groupRepo.Delete(c, id)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		c.Status(204)
		return true
	})
}
