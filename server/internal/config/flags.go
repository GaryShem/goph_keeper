package config

import (
	"github.com/spf13/viper"

	"goph_keeper/goph_server/internal/logging"
)

type ServerConfig struct {
	DbHost string
	DbPort int
	DbUser string
	DbPass string
	DbName string

	GrpcPort int
}

func ProcessFlags() (ServerConfig, error) {
	result := ServerConfig{}
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", 5432)
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "1212")
	viper.SetDefault("db_database", "goph")

	viper.SetDefault("server_port", 8080)

	if err := viper.BindEnv("db_host", "DB_HOST"); err != nil {
		return result, err
	}
	if err := viper.BindEnv("db_port", "DB_PORT"); err != nil {
		return result, err
	}
	if err := viper.BindEnv("db_user", "DB_USER"); err != nil {
		return result, err
	}
	if err := viper.BindEnv("db_password", "DB_PASSWORD"); err != nil {
		return result, err
	}
	if err := viper.BindEnv("db_database", "DB_DATABASE"); err != nil {
		return result, err
	}
	if err := viper.BindEnv("server_port", "SERVER_PORT"); err != nil {
		return result, err
	}
	result.DbHost = viper.GetString("db_host")
	result.DbPort = viper.GetInt("db_port")
	result.DbUser = viper.GetString("db_user")
	result.DbPass = viper.GetString("db_password")
	result.DbName = viper.GetString("db_database")
	result.GrpcPort = viper.GetInt("server_port")

	logging.Log().Info("db_host: ", result.DbHost)
	logging.Log().Info("db_port: ", result.DbPort)
	logging.Log().Info("db_user: ", result.DbUser)
	logging.Log().Info("db_password: ", result.DbPass)
	logging.Log().Info("db_database: ", result.DbName)

	logging.Log().Info("server_port: ", result.GrpcPort)

	return result, nil
}
