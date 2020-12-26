package main

import (
	"os"

	"github.com/useurmind/kubelab/services/projects/api/config"
	"github.com/useurmind/kubelab/services/projects/api/repository"
	"github.com/useurmind/kubelab/services/projects/api/web"

	"github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Hello projects")

	router := gin.New()

	router.Use(ginzerolog.Logger("gin"))
	router.Use(gin.Recovery())

	dbConfig := config.GetDBConfigFromEnv()
	repoFactory := repository.NewPGRepoFactory(dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.User, dbConfig.Password)
	err := repoFactory.Migrate()
	if err != nil {
		log.Error().Err(err).Msgf("Could not migrate database")
		os.Exit(1)
	}

	web.HandleGroups("groups", router, repoFactory)

	// router.POST("groups", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	log.Info().Msg("Run gin router")
	err = router.Run()
	if err != nil {
		log.Error().Err(err).Msg("Running gin failed")
		os.Exit(1)
	}
}