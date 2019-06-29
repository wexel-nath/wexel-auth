package config

import (
	"github.com/spf13/viper"
)

func Configure() {

	// Heroku Port
	viper.BindEnv("PORT")

	// Postgres URL
	viper.BindEnv("DATABASE_URL")
}

func GetPort() string {
	return viper.GetString("PORT")
}

func GetDatabaseURL() string {
	return viper.GetString("DATABASE_URL")
}
