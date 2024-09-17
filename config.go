package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey   string `json:"api_key"`
	Language string `json:"language"`
	Template string `json:"template"`
	Prompt   string `json:"prompt"`
}

func loadConfig() Config {
	configPath := getConfigPath()
	file, err := os.ReadFile(configPath)
	if err != nil {
		return Config{
			Language: "ja",
			Template: "default",
			Prompt:   getDefaultPrompt(),
		}
	}

	var config Config
	json.Unmarshal(file, &config)
	return config
}

func saveConfig(config Config) error {
	configPath := getConfigPath()

	// Ensure the directory exists
	configDir := filepath.Dir(configPath)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(configPath, file, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "gh-prai", "config.json")
}

func configureSettings(key, value string) {
	config := loadConfig()

	switch key {
	case "api_key":
		config.APIKey = value
	case "language":
		config.Language = value
	case "template":
		config.Template = value
	case "prompt":
		config.Prompt = value
	default:
		fmt.Printf("Unknown configuration key: %s\n", key)
		return
	}

	err := saveConfig(config)
	if err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
	} else {
		fmt.Printf("Configuration updated: %s\n", key)
	}
}
