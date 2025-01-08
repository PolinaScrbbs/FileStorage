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
	Secret         []byte `map structure:"secret"`
	Base_Save_Path string `map structure:"base_save_path"`
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
		Secret         string `map structure:"secret"`
		Base_Save_Path string `map structure:"base_save_path"`
	}

	if err := viper.Unmarshal(&rawConfig); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	config := &Config{
		Database:       rawConfig.Database,
		Secret:         []byte(rawConfig.Secret),
		Base_Save_Path: rawConfig.Base_Save_Path,
	}

	return config, nil
}
