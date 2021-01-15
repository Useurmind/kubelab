package main

import (
	"testing"

	"github.com/useurmind/kubelab/services/projects/api/repository"
	"github.com/useurmind/kubelab/services/projects/api/webclient"
)

type TestServiceConfig struct {
	Client *webclient.Client
}

func RunServiceUntilStopped(t *testing.T, stop chan bool) (TestServiceConfig, chan error) {
	dbSystem := repository.NewMemDBSystem()

	stopped := runWith(dbSystem, stop, true)

	return TestServiceConfig{
		Client: &webclient.Client{
			BaseUrl: "http://localhost:8080",
		},
	}, stopped
}

func EnsureServiceCorrectlyStopped(t *testing.T, stop chan bool, stopped chan error) {	
	stop <- true
	err := <- stopped
	if err != nil {
		t.Fatalf("Error returned from server %v", err)
	}
}