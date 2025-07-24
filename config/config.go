package config

import (
	"errors"
	"fmt"
	"os"
	"poebuy/utils"

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
	SoundFile string `yaml:"sound_file"`
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
		cfg.General.SoundFile = "notify.wav"
		return cfg, ErrorNoConfigFile
	}

	err := cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	cfg.General.Poesessid, _ = utils.Decrypt(cfg.General.Poesessid)

	return cfg, nil
}


// Save saves the config to file, encrypting the Poesessid
func (cfg *Config) Save() error {
    // Make a copy of cfg to avoid modifying original in-memory struct
    tempCfg := *cfg

    encPoe, err := utils.Encrypt(cfg.General.Poesessid)
    if err != nil {
        return err
    }
    tempCfg.General.Poesessid = encPoe

    err = utils.WriteStructToYAMLFile("config.yaml", &tempCfg)
    if err != nil {
        return err
    }
    return nil
}

func configFileExists() bool {
	_, err := os.Stat("config.yaml")
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
