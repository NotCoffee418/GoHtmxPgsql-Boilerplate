package common

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Error: environment variable %s not set", key)
	}
	return value
}
