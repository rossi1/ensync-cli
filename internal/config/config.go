package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Debug   bool   `mapstructure:"debug"`
}

func Load() (*Config, error) {
	config := &Config{}

	viper.SetDefault("base_url", "http://localhost:8080/api/v1/ensync")
	viper.SetDefault("debug", false)

	// Environment variables
	viper.AutomaticEnv()

	// Config file
	configDir := getConfigDir()
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate required fields
	if config.APIKey == "" {
		config.APIKey = os.Getenv("ENSYNC_API_KEY")
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	return config, nil
}

func getConfigDir() string {
	if configDir := os.Getenv("ENSYNC_CONFIG_DIR"); configDir != "" {
		return configDir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return filepath.Join(home, ".ensync")
}
