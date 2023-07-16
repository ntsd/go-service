package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

// Config is environment config
type Config struct {
	AppPort          string
	PostgresWriteURL string
	PostgresReadURL  string
	DevMode          bool
	Prefork          bool
	PrivateKeyPath   string
	PublicKeyPath    string
	HashSalt         string
}

// NewConfig returns environment config
func NewConfig() *Config {
	return &Config{
		AppPort:          defaultENV("APP_PORT", "8080"),
		PostgresWriteURL: requiredENV("POSTGRES_URL"),
		PostgresReadURL: defaultENV(
			"POSTGRES_READ_URL",
			requiredENV("POSTGRES_URL"),
		),
		DevMode:        defaultENV("DEV_MODE", "false") == "true",
		Prefork:        defaultENV("PREFORK", "true") == "true",
		PrivateKeyPath: requiredENV("ES256_PRIVATE_KEY"),
		PublicKeyPath:  requiredENV("ES256_PUBLIC_KEY"),
		HashSalt:       defaultENV("HASH_SALT", "change_this_salt"),
	}
}

// requiredENV returns environment variable value, panic if not found
func requiredENV(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal(fmt.Errorf("`%s` env is required", key))
	}
	return value
}

// defaultENV will return environment variable or default value if it's empty
func defaultENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

// strToInt convert string to int, panic if not found
func strToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(fmt.Errorf("`%s` env value can not convert to int", str))
	}
	return i
}
