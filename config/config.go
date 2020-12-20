package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var envFlag = flag.String("env", "dev", "environment, e.g. development or production")

func init() {
	flag.Parse()

	appDir, _ := findAppDir(*envFlag)
	path := joinPath(*appDir, *envFlag)

	os.Setenv("APP_ENV", *envFlag)

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
