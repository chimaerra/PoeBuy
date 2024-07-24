package config

import (
	"fmt"
	"os"
	"poebuy/utils"

	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

var ErrorNoConfigFile = errors.New("Config file not found")

// Config is the main configuration struct
type Config struct {
	General General `yaml:"general"`
	Trade   Trade   `yaml:"trade"`
}

type General struct {
	Poesessid string `yaml:"poesessid"`
}

type Trade struct {
	League string `yaml:"league"`
	Links  []Link `yaml:"links"`
}

type Link struct {
	Name    string `yaml:"name"`
	Code    string `yaml:"code"`
	Delay   int64  `yaml:"delay"`
	IsActiv bool   `yaml:"-"`
}

// LoadConfig loads the config from config file
func LoadConfig() (*Config, error) {

	cfg := &Config{}

	if !configFileExists() {
		return cfg, ErrorNoConfigFile
	}

	err := cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	cfg.General.Poesessid, _ = utils.Decrypt(cfg.General.Poesessid)

	return cfg, nil
}

func (cfg *Config) Save() {
	encPoe, _ := utils.Encrypt(cfg.General.Poesessid)
	cfg.General.Poesessid = encPoe

	utils.WriteStructToYAMLFile("config.yaml", cfg)
}

func configFileExists() bool {
	_, err := os.Stat("config.yaml")
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
