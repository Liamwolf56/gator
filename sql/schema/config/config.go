package config

import (
    "encoding/json"
    "os"
)

type Config struct {
    DBUrl          string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}

func LoadConfig() (*Config, error) {
    file, err := os.ReadFile(os.Getenv("HOME") + "/.gatorconfig.json")
    if err != nil {
        return nil, err
    }

    var cfg Config
    err = json.Unmarshal(file, &cfg)
    return &cfg, err
}

func SaveConfig(cfg *Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(os.Getenv("HOME") + "/.gatorconfig.json", data, 0644)
}
