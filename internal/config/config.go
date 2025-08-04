package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// HTTPServer holds HTTP server-related configuration.
// The field 'Addr' will be mapped from the YAML key "address".
type HTTPServer struct {
	Addr string `yaml:"address"`
}

// Config defines the overall application configuration structure.
// Fields are loaded from a YAML config file and can also be overridden via environment variables.
type Config struct {
	Env         string                                   `yaml:"env" env:"ENV" env-required:"true"` // Application environment (e.g., dev, prod)
	StoragePath string                                   `yaml:"storage_path" env-required:"true"`  // Path to the storage database
	HTTPServer  `yaml:"http_server" env-required:"true"` // Embedded HTTP server configuration
}

// MustLoad reads the configuration from file and exits the program if it fails.
// It uses cleanenv to load values from a YAML file into the Config struct.
func MustLoad() *Config {
	var cfg Config
	// Storage folder created if not exists
	if err := CreateStorageFolderIfNotExists(); err != nil {
		log.Fatalf("Failed to create storage folder: %v", err)
	}
	// Attempt to read the configuration file
	err := cleanenv.ReadConfig("config/local.yaml", &cfg)
	if err != nil {
		// Exit immediately if config fails to load
		log.Fatalf("Failed to load config: %v", err.Error())
	}

	// Return pointer to loaded config
	return &cfg
}

func CreateStorageFolderIfNotExists() error {
	_, err := os.Stat("storage")
	if os.IsNotExist(err) {
		err = os.Mkdir("storage", 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
