package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the main configuration struct
type Config struct {
	User struct {
		Poesessid string `yaml:"poesessid" env:"POESESSID"`
	} `yaml:"user"`
	Trade struct {
		Links []string `yaml:"links"`
	} `yaml:"trade"`
}

// LoadConfig loads the config from the environment variables
func LoadConfig() (*Config, error) {

	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Println("error reading environment")
		return nil, err
	}
	err = cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		log.Println("error reading config file")
		return nil, err
	}

	return cfg, nil
}
