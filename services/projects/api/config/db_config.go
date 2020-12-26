package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type DBConfig struct {
	Host string
	Port string 
	DBName string
	User string
	Password string
}

func GetDBConfigFromEnv() DBConfig {
	dbConfig := DBConfig{}
	envconfig.Process("KUBELAB_DB", &dbConfig)

	log.Info().Str("host", dbConfig.Host).Str("port", dbConfig.Port).Str("db_name", dbConfig.DBName).Str("user", dbConfig.User).Msg("Read DB configuration")

	return dbConfig
}