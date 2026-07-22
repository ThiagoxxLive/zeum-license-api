package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort         string
	DBHost          string
	DBPort          string
	DBName          string
	DBUser          string
	DBPassword      string
	DefuseSecret    string
	KeycloakBaseURL string
	KeycloakRealm   string
}

func Load() *Config {

	return &Config{
		AppPort:         getEnv("APP_PORT", "8080"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "3306"),
		DBName:          getEnv("DB_NAME", ""),
		DBUser:          getEnv("DB_USER", ""),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		DefuseSecret:    getEnv("DEFUSE_SECRET", ""),
		KeycloakBaseURL: getEnv("KEYCLOAK_BASE_URL", ""),
		KeycloakRealm:   getEnv("KEYCLOAK_REALM", ""),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
