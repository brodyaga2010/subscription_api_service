package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   Server   `yaml:"server" validate:"required"`
	Database Database `yaml:"database" validate:"required"`
	Logger   Logger   `yaml:"logger" validate:"required"`
}

func Load() (*Config, error) {
	path := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	if path == nil || *path == "" {
		return nil, fmt.Errorf("cmd flag 'config' is nil or empty")
	}

	file, err := os.Open(*path)
	if err != nil {
		return nil, fmt.Errorf("failed to opexn config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	log.Println("config: loaded successfully!")
	return &config, nil
}
