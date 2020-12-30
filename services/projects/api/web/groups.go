package web

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/useurmind/kubelab/services/projects/api/models"
	"github.com/useurmind/kubelab/services/projects/api/repository"
)

func HandleGroups(basePath string, router *gin.Engine, repoFactory repository.RepoFactory) {
	// Create a group
	router.POST(basePath, func(c *gin.Context) {
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

		if group.Id != 0 {
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
	})

	// GET group by id
	router.GET(basePath + "/:id", func(c *gin.Context) {
		repo, err := repoFactory.GetGroupRepo(c)
		if err != nil {
			AbortWithInternalError(c, err)
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			AbortWithBadRequest(c, "Could not parse input parameter id", err)
			return
		}

		group, err := repo.Get(c, id)
		if err != nil {
			AbortWithInternalError(c, err)
			return
		}

		c.JSON(200, group)
	})

	// List groups
	router.GET(basePath, func(c *gin.Context) {
		repo, err := repoFactory.GetGroupRepo(c)
		if err != nil {
			AbortWithInternalError(c, err)
			return
		}

		startIndex := 0
		if c.Param("start") != "" {
			startIndex, err = strconv.Atoi(c.Param("start"))
			if err != nil {
				AbortWithBadRequest(c, "Could not parse input parameter start", err)
				return
			}
		} 
		count := 50
		if c.Param("count") != "" {
			count, err = strconv.Atoi(c.Param("count"))
			if err != nil {
				AbortWithBadRequest(c, "Could not parse input parameter count", err)
				return
			}
		} 

		groups, err := repo.List(c, startIndex, count)
		if err != nil {
			AbortWithInternalError(c, err)
			return
		}

		c.JSON(200, groups)
	})
}