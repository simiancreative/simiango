package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var envFlag string

func init() {
	flag.StringVar(
		&envFlag,
		"env",
		"",
		"environment, e.g. development or production",
	)

	flag.Parse()

	if len(envFlag) == 0 {
		envFlag = os.Getenv("APP_ENV")
	}

	if len(envFlag) == 0 {
		envFlag = "dev"
	}

	os.Setenv("APP_ENV", envFlag)

	appDir, _ := findAppDir(envFlag)
	path := joinPath(*appDir, envFlag)

	_ = godotenv.Load(path)
}

func findAppDir(env string) (*string, error) {
	appDir, err := filepath.Abs("")
	if err != nil {
		return nil, err
	}

	saved := appDir

	for i := 0; i < 10; i++ {
		cfgPath := joinPath(appDir, env)
		_, err := os.Stat(cfgPath)
		if err == nil {
			return &appDir, nil
		}

		if appDir == "." {
			break
		}

		appDir = filepath.Dir(appDir)
	}

	return &saved, nil
}

func joinPath(appDir, env string) string {
	return filepath.Join(appDir, "config", fmt.Sprintf(".env.%s.local", env))
}
