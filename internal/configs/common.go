package configs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nekomeowww/xo"
	"github.com/spf13/viper"
)

const ConfigFilePathEnvName = "CONFIG_FILE_PATH"

func getConfigFilePath(configFilePath string) string {
	if configFilePath != "" {
		return configFilePath
	}

	envPath := os.Getenv(ConfigFilePathEnvName)
	if envPath != "" {
		return envPath
	}

	configPath := xo.RelativePathBasedOnPwdOf("config/config.yaml")

	return configPath
}

var (
	possibleConfigPathsForUnitTest = []string{
		"config.local.yml",
		"config.local.yaml",
		"config.test.yml",
		"config.test.yaml",
		"config.example.yml",
		"config.example.yaml",
	}
)

func tryToMatchConfigPathForUnitTest(configFilePath string) string {
	if getConfigFilePath(configFilePath) != "" {
		return configFilePath
	}

	for _, path := range possibleConfigPathsForUnitTest {
		stat, err := os.Stat(filepath.Join(xo.RelativePathOf("../../config"), path))
		if err == nil {
			if stat.IsDir() {
				panic("config file path is a directory: " + path)
			}

			return path
		}
		if errors.Is(err, os.ErrNotExist) {
			continue
		} else {
			panic(err)
		}
	}

	return ""
}

func loadEnvConfig(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := godotenv.Load(xo.RelativePathBasedOnPwdOf("./.env"))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}

			return nil
		}

		return err
	}

	return nil
}

func readConfig(path string) error {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { //nolint:errorlint
			return nil
		}
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("error occurred when read in config, error is: %s, err: %w", fmt.Sprintf("%T", err), err)
	}

	return nil
}
