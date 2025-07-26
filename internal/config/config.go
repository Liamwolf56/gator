package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Config represents the structure of the JSON configuration file.
type Config struct {
	DBURL          string `json:"db_url"`
	CurrentUserName string `json:"current_user_name,omitempty"` // omitempty to not include if empty
}

// getConfigFilePath returns the full path to the configuration file.
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

// Read reads the JSON config file from the HOME directory and returns a Config struct.
func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// write writes the given Config struct to the JSON file.
func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Use MarshalIndent for prettified JSON, as often preferred for config files
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Ensure the file permissions are appropriate (e.g., 0644 for rw-r--r--)
	return ioutil.WriteFile(filePath, data, 0644)
}

// SetUser sets the current_user_name field and writes the updated config to disk.
func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}
