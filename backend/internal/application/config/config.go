package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server         ServerConfig         `yaml:"server"`
	OpenRouter     OpenRouterConfig     `yaml:"openrouter"`
	Models         ModelsConfig         `yaml:"models"`
	Categorization CategorizationConfig `yaml:"categorization"`
}

type ServerConfig struct {
	Port         int           `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type OpenRouterConfig struct {
	APIKey  string        `yaml:"-"` // Loaded from env
	BaseURL string        `yaml:"base_url"`
	Timeout time.Duration `yaml:"timeout"`
}

type ModelsConfig struct {
	Categorizer string `yaml:"categorizer"`
	Easy        string `yaml:"easy"`
	Advanced    string `yaml:"advanced"`
	Coding      string `yaml:"coding"`
	Image       string `yaml:"image"`
	ImageHard   string `yaml:"image_hard"`
}

type CategorizationConfig struct {
	Prompt string `yaml:"prompt"`
}

func Load(configPath string) (*Config, error) {
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Load sensitive data from environment
	cfg.OpenRouter.APIKey = os.Getenv("OPENROUTER_API_KEY")
	if cfg.OpenRouter.APIKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY environment variable is required")
	}

	// Override base URL if set in env
	if baseURL := os.Getenv("OPENROUTER_BASE_URL"); baseURL != "" {
		cfg.OpenRouter.BaseURL = baseURL
	}

	// Override port if set in env
	if port := os.Getenv("SERVER_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			cfg.Server.Port = p
		}
	}

	return &cfg, nil
}

func (c *Config) GetModelForCategory(category string) (string, error) {
	switch category {
	case "easy":
		return c.Models.Easy, nil
	case "advanced":
		return c.Models.Advanced, nil
	case "coding":
		return c.Models.Coding, nil
	case "image":
		return c.Models.Image, nil
	default:
		return "", fmt.Errorf("unknown category: %s", category)
	}
}
