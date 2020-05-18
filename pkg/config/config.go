package config

import (
	"time"

	"github.com/spf13/viper"
)

func Configure() {

	// Auth config
	viper.Set("JWT_ISSUER", "wexel-auth.herokuapp.com")
	viper.Set("JWT_EXPIRY", 5 * 60)
	viper.Set("SESSION_EXPIRY", 60 * time.Minute)
	viper.Set("SERVICE_NAME", "authentication")
	_ = viper.BindEnv("PRIVATE_KEY_PATH")
	_ = viper.BindEnv("PUBLIC_KEY_PATH")

	// Heroku Port
	viper.SetDefault("PORT", "4000")
	_ = viper.BindEnv("PORT")

	// Postgres URL
	viper.SetDefault("DATABASE_URL", "")
	_ = viper.BindEnv("DATABASE_URL")
}

func GetJwtIssuer() string {
	return viper.GetString("JWT_ISSUER")
}

func GetJwtExpiry() int64 {
	return viper.GetInt64("JWT_EXPIRY")
}

func GetPublicKeyPath() string {
	return viper.GetString("PUBLIC_KEY_PATH")
}

func GetPrivateKeyPath() string {
	return viper.GetString("PRIVATE_KEY_PATH")
}

func GetSessionExpiry() time.Duration {
	return viper.Get("SESSION_EXPIRY").(time.Duration)
}

func GetServiceName() string {
	return viper.GetString("SERVICE_NAME")
}

func GetPort() string {
	return viper.GetString("PORT")
}

func GetDatabaseURL() string {
	return viper.GetString("DATABASE_URL")
}
