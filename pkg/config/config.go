package config

import (
	"github.com/spf13/viper"
)

func Configure() {

	// Auth config
	viper.Set("JWT_ISSUER", "wexel-auth.herokuapp.com")
	viper.Set("JWT_EXPIRY", 5 * 60)
	viper.Set("SESSION_EXPIRY", 60 * 60)

	// Heroku Port
	viper.BindEnv("PORT")

	// Postgres URL
	viper.BindEnv("DATABASE_URL")
}

func GetJwtIssuer() string {
	return viper.GetString("JWT_ISSUER")
}

func GetJwtExpiry() int64 {
	return viper.GetInt64("JWT_EXPIRY")
}

func GetSessionExpiry() int64 {
	return viper.GetInt64("SESSION_EXPIRY")
}

func GetPort() string {
	return viper.GetString("PORT")
}

func GetDatabaseURL() string {
	return viper.GetString("DATABASE_URL")
}
