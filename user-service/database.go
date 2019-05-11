package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

func CreateConnection() (*gorm.DB, error) {
	// postgras db connection vars
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "")
	dbName := getEnv("DB_NAME", "")
	password := getEnv("DB_PASSWORD", "")

	return gorm.Open(
		"postgras",
		fmt.Sprintf(
			"postgras://%s:%s@%s/%s?sslmode=disable",
			user, password, host, dbName,
		),
	)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
