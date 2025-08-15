package config

type Logger struct {
	FilePath string `yaml:"file_path" validate:"required"`
}
