package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configPath)
	if err != nil {
    return Config{}, fmt.Errorf("issue opening file: %v", err)
	}
	defer file.Close()

  byteVal, err := io.ReadAll(file)
  if err != nil {
    return Config{}, fmt.Errorf("issue reading file: %v", err)
  }

	var config Config
	err = json.Unmarshal(byteVal, &config)
	if err != nil {
    return config, fmt.Errorf("issue decoding: %v", err)
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username

	configPath, err := getConfigFilePath()
	if err != nil {
    fmt.Printf("Issue getting config path: %v", err)
		return err
	}
	jsonData, err := json.MarshalIndent(c, "", "  ")
  if err != nil {
    fmt.Printf("Issue marshaling data: %v", err)
    return err
  }

  err = os.WriteFile(configPath, jsonData, 0644)
  if err != nil {
    fmt.Printf("Issue writing file: %v", err)
    return err
  }

  fmt.Println("json data written to ", configPath)
	return nil
}

func getConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return userHomeDir + configFileName, nil
}
