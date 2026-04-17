// Package config reads the config files
package config

import (
	"encoding/json"
	"os"
)

const configFile = "/.gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return homeDir + configFile
}

func Read() (Config, error) {
	content, err := os.ReadFile(getConfigPath())
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return config, err
	}

	return config, nil
}

func (conf *Config) SetUser() error {
	jsonData, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = os.WriteFile(getConfigPath(), jsonData, 0o644)
	if err != nil {
		return err
	}

	return nil
}
