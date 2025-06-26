package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configJson = "/.gatorconfig.json"

type Config struct {
	DbUrl           *string `json:"db_url"`
	CurrentUserName *string `json:"current_user_name"`
}

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmtErr := fmt.Errorf("Error getting home dir: %v", err)
		return Config{}, fmtErr
	}

	home += configJson

	confFile, err := os.Open(home)
	if err != nil {
		fmtErr := fmt.Errorf("Error opening config file: %v", err)
		return Config{}, fmtErr
	}
	defer confFile.Close()

	data, err := io.ReadAll(confFile)
	if err != nil {
		fmtErr := fmt.Errorf("Error decoding config file: %v", err)
		return Config{}, fmtErr
	}

	conf := Config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		fmtErr := fmt.Errorf("Error unmarshalling config data: %v", err)
		return Config{}, fmtErr
	}

	return conf, nil
}

func (c Config) SetUser(username string) error {
	c.CurrentUserName = &username

	home, err := os.UserHomeDir()
	if err != nil {
		fmtErr := fmt.Errorf("Error getting home dir: %v", err)
		return fmtErr
	}

	home += configJson

	data, err := json.Marshal(c)
	if err != nil {
		fmtErr := fmt.Errorf("Error marshalling config: %v", err)
		return fmtErr
	}

	if err := os.WriteFile(home, data, 0644); err != nil {
		fmtErr := fmt.Errorf("Error writing data to file: %v", err)
		return fmtErr
	}

	return nil
}
