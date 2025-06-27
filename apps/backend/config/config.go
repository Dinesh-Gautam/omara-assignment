package config

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	FrontendURL                  string
	PostgresHost                 string
	PostgresPort                 string
	PostgresUser                 string
	PostgresPassword             string
	PostgresDB                   string
	GoogleApplicationCredentials string
	GeminiAPIKey                 string
	GCSBucketName                string
}

// AppConfig is a global variable that holds the application configuration
var AppConfig *Config

// LoadConfig loads configuration from a .env file and environment variables
func LoadConfig() error {
	// We assume the .env file is in the same directory as the executable
	// or in the project root when running with `go run`.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Reading configuration from environment variables.")
	}

	AppConfig = &Config{
		FrontendURL:                  os.Getenv("FRONTEND_URL"),
		PostgresHost:                 os.Getenv("POSTGRES_HOST"),
		PostgresPort:                 os.Getenv("POSTGRES_PORT"),
		PostgresUser:                 os.Getenv("POSTGRES_USER"),
		PostgresPassword:             os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:                   os.Getenv("POSTGRES_DB"),
		GoogleApplicationCredentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		GeminiAPIKey:                 os.Getenv("GEMINI_API_KEY"),
		GCSBucketName:                os.Getenv("GCS_BUCKET_NAME"),
	}

	// Validate that all required environment variables are set.
	// Validate required fields
	requiredVars := map[string]string{
		"FRONTEND_URL":    AppConfig.FrontendURL,
		"POSTGRES_HOST":   AppConfig.PostgresHost,
		"POSTGRES_PORT":   AppConfig.PostgresPort,
		"POSTGRES_USER":   AppConfig.PostgresUser,
		"POSTGRES_DB":     AppConfig.PostgresDB,
		"GEMINI_API_KEY":  AppConfig.GeminiAPIKey,
		"GCS_BUCKET_NAME": AppConfig.GCSBucketName,
	}

	var missingVars []string
	for key, value := range requiredVars {
		if value == "" {
			missingVars = append(missingVars, key)
		}
	}

	// Also check for PostgresPassword separately to avoid logging it
	if AppConfig.PostgresPassword == "" {
		missingVars = append(missingVars, "POSTGRES_PASSWORD")
	}

	if len(missingVars) > 0 {
		sort.Strings(missingVars) // Sort for consistent error messages
		return fmt.Errorf("FATAL: missing required environment variables: %s", strings.Join(missingVars, ", "))
	}
	// Note: GOOGLE_APPLICATION_CREDENTIALS is considered optional, as the Firebase SDK
	// can use default credentials in a Google Cloud environment.

	return nil
}

// GetDBDSN returns the full database connection string.
func (c *Config) GetDBDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
}
