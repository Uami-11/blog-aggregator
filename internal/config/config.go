// Package config
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	var readConfig Config

	filePath, err := getConfigFilePath()
	if err != nil {
		return readConfig
	}

	file, err := os.Open(filePath)
	if err != nil {
		return readConfig
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := json.NewDecoder(file)

	_ = decoder.Decode(&readConfig)

	return readConfig
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	err := write(*cfg)

	return err
}

func getConfigFilePath() (string, error) {
	filePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath = filepath.Join(filePath, configFileName)

	return filePath, nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, _ := os.Create(filePath)

	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(cfg)

	return err
}
