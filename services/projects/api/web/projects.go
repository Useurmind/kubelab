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

func HandleProjects(basePath string, router *gin.Engine, dbSystem repository.DBSystem) {
	// Create a project
	router.POST(basePath, func(c *gin.Context) { UseDBContext(c, dbSystem, createProject) })

	// Update a project
	router.PUT(basePath+"/:id", func(c *gin.Context) { UseDBContext(c, dbSystem, updateProject) })

	// GET project by id
	router.GET(basePath+"/:id", func(c *gin.Context) { UseDBContext(c, dbSystem, getProjectByID) })

	// List project
	router.GET(basePath, func(c *gin.Context) { UseDBContext(c, dbSystem, listProjects) })

	// DELETE project
	router.DELETE(basePath+"/:id", func(c *gin.Context) { UseDBContext(c, dbSystem, deleteProject) })
}

func createProject(c *gin.Context, dbContext repository.DBContext) {
	projectRepo, err := dbContext.GetProjectRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}
	groupRepo, err := dbContext.GetGroupRepo(c)
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
	project := models.Project{}
	err = json.Unmarshal(jsonBody, &project)
	if err != nil {
		log.Error().Str("jsonbody", string(jsonBody)).Msg("Could not parse json body")
		AbortWithBadRequest(c, "Could not parse body as project object.", err)
		return
	}

	if !project.IsNew() {
		AbortWithBadRequest(c, "Project already has an id, did you intend to update an existing project? Use PUT /projects/:id instead.", nil)
		return
	}

	if project.GroupId == 0 {
		AbortWithBadRequest(c, "Project has no group id set, which is required.", nil)
		return
	}

	CommitOnSuccess(c, dbContext, func() bool {
		projectResult, err := projectRepo.CreateOrUpdate(c, &project)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		group, err := groupRepo.Get(c, project.GroupId)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		err = group.InsertProjectRef(projectResult)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		_, err = groupRepo.CreateOrUpdate(c, group)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		c.JSON(200, projectResult)
		return true
	})
}

func updateProject(c *gin.Context, dbContext repository.DBContext) {
	projectRepo, err := dbContext.GetProjectRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}
	groupRepo, err := dbContext.GetGroupRepo(c)
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

	project := models.Project{}
	err = json.Unmarshal(jsonBody, &project)
	if err != nil {
		log.Error().Str("jsonbody", string(jsonBody)).Msg("Could not parse json body")
		AbortWithBadRequest(c, "Could not parse body as project object.", err)
		return
	}

	if project.IsNew() {
		AbortWithBadRequest(c, fmt.Sprintf("Project id (%d) is zero or not set (path id %d). This indicates a new group. Use POST /projects instead.", project.Id, id), nil)
		return
	}

	if project.Id != id {
		AbortWithBadRequest(c, fmt.Sprintf("Project id %d does not match id %d in url, did you really want to put this project?", project.Id, id), nil)
		return
	}

	CommitOnSuccess(c, dbContext, func() bool {
		projectResult, err := projectRepo.CreateOrUpdate(c, &project)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		group, err := groupRepo.Get(c, project.GroupId)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		err = group.UpdateProjectRef(projectResult)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		_, err = groupRepo.CreateOrUpdate(c, group)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		c.JSON(200, projectResult)
		return true
	})
}

func getProjectByID(c *gin.Context, dbContext repository.DBContext) {
	repo, err := dbContext.GetProjectRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	id, err := getIntParam(c, "id")
	if err != nil {
		return
	}

	project, err := repo.Get(c, id)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	if project == nil {
		AbortWithNotFoundError(c, nil)
	}

	c.JSON(200, project)
}

func deleteProject(c *gin.Context, dbContext repository.DBContext) {
	projectRepo, err := dbContext.GetProjectRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}
	groupRepo, err := dbContext.GetGroupRepo(c)
	if err != nil {
		AbortWithInternalError(c, err)
		return
	}

	id, err := getIntParam(c, "id")
	if err != nil {
		return
	}

	CommitOnSuccess(c, dbContext, func() bool {
		project, err := projectRepo.Get(c, id)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		err = projectRepo.Delete(c, id)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		group, err := groupRepo.Get(c, project.GroupId)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		err = group.DeleteProjectRef(project)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		_, err = groupRepo.CreateOrUpdate(c, group)
		if err != nil {
			AbortWithInternalError(c, err)
			return false
		}

		c.Status(204)
		return true
	})
}
