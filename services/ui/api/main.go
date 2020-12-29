package main

import (
	"os"

	"github.com/Useurmind/spas/handler"
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/useurmind/kubelab/services/ui/api/models"
)

func getConfig() models.UIConfig {
	uiConfig := models.UIConfig{}
	envconfig.Process("KUBELAB_UI", &uiConfig)

	log.Info().Str("uiBaseUrl", uiConfig.UiBaseUrl).Str("projectsBaseUrl", uiConfig.ProjectsBaseUrl).Msg("Read UI configuration")

	return uiConfig
}

func main() {
	log.Info().Msg("Starting kubelab ui")

	router := gin.New()

	router.Use(ginzerolog.Logger("gin"))
	router.Use(gin.Recovery())

	spasOptions := handler.Options{
		// only these two options are actually used by the handler
		// the rest is only used to start the spa server and not relevant here
		ServeFolder:   "www",
		HTMLIndexFile: "index.html",
	}
	spasHandler := handler.NewSPASHandler(&spasOptions)

	uiConfig := getConfig()

	router.Handle("GET", "/ui/*", gin.WrapH(spasHandler))
	router.Handle("GET", "/config", func(c *gin.Context) {
		c.JSON(200, uiConfig)
	})

	log.Info().Msg("Run gin router")
	err := router.Run()
	if err != nil {
		log.Error().Err(err).Msg("Running gin failed")
		os.Exit(1)
	}
}
