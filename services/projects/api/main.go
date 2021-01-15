package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/useurmind/kubelab/services/projects/api/config"
	"github.com/useurmind/kubelab/services/projects/api/repository"
	"github.com/useurmind/kubelab/services/projects/api/web"

	"github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting kubelab projects")

	dbConfig := config.GetDBConfigFromEnv()
	dbSystem := repository.NewPGDBSystem(dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.User, dbConfig.Password)
	
	stop := make(chan bool)
	stopped := runWith(dbSystem, stop, false)
	err := <- stopped
	if err != nil {
		log.Error().Err(err).Msg("Starting kubelab projects failed")
		os.Exit(1)
	}
}

func runWith(dbSystem repository.DBSystem, stop chan bool, local bool) chan error {
	err := migrateDB(dbSystem)
	if err != nil {
		stopped := make(chan error)
		stopped <- err
		return stopped
	}

	router := gin.New()

	router.Use(ginzerolog.Logger("gin"))
	router.Use(cors.Default())
	router.Use(gin.Recovery())

	web.HandleGroups("groups", router, dbSystem)

	// router.POST("groups", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	stopped := runGinServer(stop, router, local)
	return stopped
}

func migrateDB(dbSystem repository.DBSystem) error {
	dbContext := dbSystem.NewContext()
	defer dbContext.Close()
	err := dbContext.Migrate()
	if err != nil {
		err2 := fmt.Errorf("Could not migrated database: %w", err)
		return err2
	}

	return nil
}

func runGinServer(stop chan bool, router *gin.Engine, local bool) chan error {
	interf := ""
	if local {
		interf = "localhost"
	}
	port := "8080"


	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", interf, port),
		Handler: router,
	}
	stopped := make(chan error)

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Error during ListenAndServe")
			stopped <- err
		}
	}()

	go func() {
		// detect whether stop was requested
		<-stop
		log.Info().Msg("Shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
			stopped <- nil
		}()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server Shutdown")
			stopped <- err
		}
		log.Info().Msg("Server exiting")
	}()

	return stopped
}
