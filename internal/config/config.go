package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

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

	home += "/.gatorconfig.json"

	confFile, err := os.Open(home)
	if err != nil {
		fmtErr := fmt.Errorf("Error reading config file: %v", err)
		return Config{}, fmtErr
	}
	defer confFile.Close()

	data, err := io.ReadAll(confFile)
	if err != nil {
		fmtErr := fmt.Errorf("Error decoding config file: %v", err)
		return Config{}, fmtErr
	}
	fmt.Println(fmt.Sprintf("%v", string(data)))

	conf := Config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		fmtErr := fmt.Errorf("Error unmarshalling config data: %v", err)
		return Config{}, fmtErr
	}

	return conf, nil
}

func (Config) SetUser(username string) error {
	return nil
}
