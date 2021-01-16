package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/useurmind/kubelab/services/projects/api/models"
	"github.com/useurmind/kubelab/services/projects/api/webclient"
)

func TestGroupCanAddGetAndDelete(t *testing.T) {
	stop := make(chan bool)
	config, stopped := RunServiceUntilStopped(t, stop)

	group := models.Group{
		Name: "group1",
	}

	createdGroup, err := config.Client.Groups().Create(&group)
	if err != nil {
		t.Fatalf("Error while creating group %v", err)
	}

	getGroup, err := config.Client.Groups().Get(createdGroup.Id)
	if err != nil {
		t.Fatalf("Error while getting group %v", err)
	}

	err = config.Client.Groups().Delete(createdGroup.Id)
	if err != nil {
		t.Fatalf("Error while deleting group %v", err)
	}

	_, err = config.Client.Groups().Get(createdGroup.Id)
	if err != webclient.ErrorNotFound {
		t.Errorf("Expected 404 error for deleted group but got %v", err)
	}

	assert.Greater(t, createdGroup.Id, int64(0), "The created group should have an id set greater zero")
	assert.Equal(t, createdGroup.Id, getGroup.Id, "Both groups should have the same id")
	assert.Equal(t, createdGroup.Name, getGroup.Name, "Both groups should have the same name")

	EnsureServiceCorrectlyStopped(t, stop, stopped)
}

func TestSubgroupsGetIdWhenCreating(t *testing.T) {
	stop := make(chan bool)
	config, stopped := RunServiceUntilStopped(t, stop)

	group := models.Group{
		Name: "group1",
		Subgroups: []*models.Group{
			{
				Name: "subgroup1",
			},
			{
				Name: "subgroup2",
			},
		},
	}

	createdGroup, err := config.Client.Groups().Create(&group)
	if err != nil {
		t.Fatalf("Error while creating group %v", err)
	}

	
	assert.Greater(t, createdGroup.Id, int64(0), "The created group should have an id set greater zero")
	assert.Greater(t, createdGroup.Subgroups[0].Id, int64(0), "The created subgroup 1 should have an id set greater zero")
	assert.Greater(t, createdGroup.Subgroups[1].Id, int64(0), "The created subgroup 2 should have an id set greater zero")

	EnsureServiceCorrectlyStopped(t, stop, stopped)
}

func TestListingGroups(t *testing.T) {
	stop := make(chan bool)
	config, stopped := RunServiceUntilStopped(t, stop)

	group1 := models.Group{
		Name: "group1",
	}
	group2 := models.Group{
		Name: "group2",
	}

	createdGroup1, err := config.Client.Groups().Create(&group1)
	if err != nil {
		t.Fatalf("Error while creating group 1 %v", err)
	}
	createdGroup2, err := config.Client.Groups().Create(&group2)
	if err != nil {
		t.Fatalf("Error while creating group 2 %v", err)
	}

	createdGroups := map[int64]*models.Group{
		createdGroup1.Id: createdGroup1,
		createdGroup2.Id: createdGroup2,
	}

	groups, err := config.Client.Groups().List(0, 100)
	if err != nil {
		t.Fatalf("Error while creating group %v", err)
	}

	assert.Equal(t, len(groups), 2, "There should be 2 groups in the list but got %d", len(groups))
	for _, group := range groups {
		createdGroup := createdGroups[group.Id]
		assert.Equal(t, group.Name, createdGroup.Name, "Groups with id %d should have same name but '%s' != '%s'", group.Name, createdGroup.Name)
	}

	EnsureServiceCorrectlyStopped(t, stop, stopped)
}