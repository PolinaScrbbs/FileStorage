package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string
		Port     int32
		User     string
		Password string
		Name     string
	}
	Secret []byte `map structure:"secret"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var rawConfig struct {
		Database struct {
			Host     string
			Port     int32
			User     string
			Password string
			Name     string
		}
		Secret string `map structure:"secret"`
	}

	if err := viper.Unmarshal(&rawConfig); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	config := &Config{
		Database: rawConfig.Database,
		Secret:   []byte(rawConfig.Secret),
	}

	return config, nil
}
