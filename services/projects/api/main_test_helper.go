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

	address := "localhost:8081"

	stopped := runWith(dbSystem, stop, address)

	return TestServiceConfig{
		Client: &webclient.Client{
			BaseUrl: "http://" + address,
		},
	}, stopped
}

func EnsureServiceCorrectlyStopped(t *testing.T, stop chan bool, stopped chan error) {	
	stop <- true
	err := <- stopped
	if err != nil {
		t.Fatalf("Error returned from server: %v", err)
	}
}