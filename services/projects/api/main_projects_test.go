package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/useurmind/kubelab/services/projects/api/models"
	"github.com/useurmind/kubelab/services/projects/api/webclient"
)

func TestProjectsCanAddGetAndDelete(t *testing.T) {
	stop := make(chan bool)
	config, stopped := RunServiceUntilStopped(t, stop)

	group := models.Group{
		Name: "group1",
		Subgroups: []*models.Group{
			{ Name: "subgroup1" },
		},
	}
	project := models.Project{
		Name: "project1",
	}

	createdGroup, err := config.Client.Groups().Create(&group)
	if err != nil {
		t.Fatalf("Error while creating group %v", err)
	}

	project.GroupId = createdGroup.Id
	project.AssignedGroupId = createdGroup.Subgroups[0].Id

	createdProject, err := config.Client.Projects().Create(&project)
	if err != nil {
		t.Fatalf("Error while creating project %v", err)
	}

	getProject, err := config.Client.Projects().Get(createdProject.Id)
	if err != nil {
		t.Fatalf("Error while getting project %v", err)
	}

	err = config.Client.Projects().Delete(createdProject.Id)
	if err != nil {
		t.Fatalf("Error while deleting project %v", err)
	}

	_, err = config.Client.Projects().Get(createdProject.Id)
	if err != webclient.ErrorNotFound {
		t.Errorf("Expected 404 error for deleted project but got %v", err)
	}

	assert.Greater(t, createdProject.Id, int64(0), "The created proejct should have an id set greater zero")
	assert.Equal(t, createdProject.Id, getProject.Id, "Both projects should have the same id")
	assert.Equal(t, createdProject.Name, getProject.Name, "Both projects should have the same name")

	EnsureServiceCorrectlyStopped(t, stop, stopped)
}

func TestProjectsRequireGroupIdSet(t *testing.T) {
	stop := make(chan bool)
	config, stopped := RunServiceUntilStopped(t, stop)

	project := models.Project{
		Name: "project1",
	}

	_, err := config.Client.Projects().Create(&project)
	if err.Error() != "Expected response code 200, but got 400: {\"error\":\"Project has no group id set, which is required.\"}" {
		t.Fatalf("Expected validation error but got: %v", err)
	}

	EnsureServiceCorrectlyStopped(t, stop, stopped)
}

func TestProjectsRequireSubGroupIdSet(t *testing.T) {
	stop := make(chan bool)
	config, stopped := RunServiceUntilStopped(t, stop)

	group := models.Group{
		Name: "group1",
	}
	project := models.Project{
		Name: "project1",
	}

	createdGroup, err := config.Client.Groups().Create(&group)
	if err != nil {
		t.Fatalf("Error while creating group %v", err)
	}

	project.GroupId = createdGroup.Id

	_, err = config.Client.Projects().Create(&project)
	if err.Error() != "Expected response code 200, but got 400: {\"error\":\"Project has no assigned group id set, which is required.\"}" {
		t.Fatalf("Expected validation error but got: %v", err)
	}

	EnsureServiceCorrectlyStopped(t, stop, stopped)
}

func TestGroupReferencesProjectAfterCreation(t *testing.T) {
	stop := make(chan bool)
	config, stopped := RunServiceUntilStopped(t, stop)

	group, project := createGroupWithProject(t, config)
	updatedGroup, err := config.Client.Groups().Get(group.Id)
	if err != nil {
		t.Fatalf("Error while getting group %v", err)
	}

	if len(updatedGroup.Subgroups[0].Projects) != 1 {
		t.Fatalf("Expected a project in subgroup but got %d", len(group.Subgroups[0].Projects))
	}

	projectRef := updatedGroup.Subgroups[0].Projects[0]
	assertProjectRefMatch(t, project, projectRef)

	EnsureServiceCorrectlyStopped(t, stop, stopped)
}

func assertProjectRefMatch(t *testing.T, project *models.Project, projectRef *models.ProjectRef) {
	assert.Equal(t, project.Id, projectRef.Id)
	assert.Equal(t, project.Name, projectRef.Name)
	assert.Equal(t, project.Slug, projectRef.Slug)
}

func createGroupWithProject(t *testing.T, config TestServiceConfig) (*models.Group, *models.Project) {
	group := models.Group{
		Name: "group1",
		Subgroups: []*models.Group{
			{ Name: "subgroup1" },
		},
	}
	project := models.Project{
		Name: "project1",
	}

	createdGroup, err := config.Client.Groups().Create(&group)
	if err != nil {
		t.Fatalf("Error while creating group %v", err)
	}

	project.GroupId = createdGroup.Id
	project.AssignedGroupId = createdGroup.Subgroups[0].Id

	createdProject, err := config.Client.Projects().Create(&project)
	if err != nil {
		t.Fatalf("Error while creating project %v", err)
	}

	return createdGroup, createdProject
}