package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/useurmind/kubelab/services/projects/api/models"
)


func TestGroupCanAddAndGet(t *testing.T) {
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

	assert.Equal(t, createdGroup.Id, getGroup.Id, "Both groups should have the same id")
	assert.Equal(t, createdGroup.Name, getGroup.Name, "Both groups should have the same name")

	EnsureServiceCorrectlyStopped(t, stop, stopped)
}