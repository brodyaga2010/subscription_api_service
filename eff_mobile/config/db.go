package config

type Database struct {
	ConnectionString string `yaml:"connection-string" validate:"required"`
	Timeout          int    `yaml:"timeout" validate:"required"`
}
