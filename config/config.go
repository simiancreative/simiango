package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/simiancreative/simiango/logger"
)

// WithFlag is a helper function to load the environment from a flag
//
// Usage: go run main.go --dot-env=development
//
// This will load the .env.development file or .env.development.local file in
// the nearest confg directory
func WithFlag() {
	var envFlag string

	flag.StringVar(
		&envFlag,
		"dot-env",
		"",
		"environment, e.g. development or production",
	)

	flag.Parse()

	if envFlag == "" {
		logger.New().Error("No environment specified")
		return
	}

	os.Setenv("DOT_ENV", envFlag)

	Enable()
}

// Enable loads the environment from the DOT_ENV environment variable
//
// This will load the .env.development file or .env.development.local file in
// the nearest confg directory
func Enable() {
	env := os.Getenv("DOT_ENV")
	if env == "" {
		return
	}

	path, err := findConfigWithFallback(env)
	if err != nil {
		logger.New().Errorf("Failed to load config (%v) - %v", err, env)
		return
	}

	err = godotenv.Load(path)
	if err != nil {
		logger.New().Errorf("Failed to load godotenv (%v) - %v", err, path)
	}
}

func SetupTest() {
	path, err := findConfigWithFallback("test")
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(path)
	if err != nil {
		panic(err)
	}

	logger.Enable()
}

func findConfigWithFallback(env string) (string, error) {
	var path string

	baseFilePath := "config/.env." + env
	err := findConfig(baseFilePath, &path)
	if err != nil {
		logger.New().Errorf("Base path not found (%v) - %v", err, baseFilePath)
	}

	if err == nil {
		return path, nil
	}

	localFilePath := "config/.env." + env + ".local"
	err = findConfig(localFilePath, &path)
	if err != nil {
		logger.New().Errorf("Local path not found (%v) - %v", err, localFilePath)
	}

	return path, err
}

func findConfig(fileName string, dest *string) error {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Start from the current directory and move upwards to find the top-level config directory
	for {
		configFilePath := filepath.Join(currentDir, fileName)

		// Check if the file exists
		if _, err := os.Stat(configFilePath); err == nil {
			*dest = configFilePath
			return nil // Found the file, return its path
		}

		// Move up one directory level
		parentDir := filepath.Dir(currentDir)

		// Check if we've reached the root directory
		if parentDir == currentDir {
			break
		}

		currentDir = parentDir // Move up to the parent directory
	}

	// File not found in any parent directory
	return fmt.Errorf("not found in any parent directory")
}
