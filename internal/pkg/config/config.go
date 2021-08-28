package config

import (
	"os"
	"strconv"
)

type Configuration struct {
	Port     string
	StoreTtl int
}

func SetUpConfiguration() *Configuration {
	config := &Configuration{
		Port:     getEnvStr("HISTORY_SERVER_LISTEN_ADDR", "8080"),
		StoreTtl: getEnvInt("LOCATION_HISTORY_TTL_SECONDS", 60),
	}
	return config
}

func getEnvStr(variable string, defaultValue string) string {
	if env, ok := os.LookupEnv(variable); ok {
		return env
	}
	return defaultValue
}

func getEnvInt(variable string, defaultValue int) int {
	valStr, ok := os.LookupEnv(variable)
	if !ok {
		return defaultValue
	}
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}

	return defaultValue
}
