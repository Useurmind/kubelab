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

func createGroup(c *gin.Context, repoFactory repository.DBContext) {
	repo, err := repoFactory.GetGroupRepo(c)
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

	groupResult, err := repo.CreateOrUpdate(c, &group)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	c.JSON(200, groupResult)
}

func updateGroup(c *gin.Context, repoFactory repository.DBContext) {
	repo, err := repoFactory.GetGroupRepo(c)
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

	groupResult, err := repo.CreateOrUpdate(c, &group)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	c.JSON(200, groupResult)
}

func getGroupByID(c *gin.Context, repoFactory repository.DBContext) {
	repo, err := repoFactory.GetGroupRepo(c)
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

	c.JSON(200, group)
}

func listGroups(c *gin.Context, repoFactory repository.DBContext) {
	repo, err := repoFactory.GetGroupRepo(c)
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

func deleteGroup(c *gin.Context, repoFactory repository.DBContext) {
	repo, err := repoFactory.GetGroupRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	id, err := getIntParam(c, "id")
	if err != nil {
		return
	}

	err = repo.Delete(c, id)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	c.Status(204)
}
