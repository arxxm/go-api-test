package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	PostgresHost string
	PostgresPort int
	PostgresName string
	PostgresUser string
	PostgresPass string
}

func New() (*Config, error) {
	absPath, err := os.Getwd()
	if err != nil {

		return nil, fmt.Errorf("error getting absolute path: %s", err)
	}
	viper.SetConfigFile(absPath + "/.env")

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %s", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling configuration: %s", err)
	}

	log.Printf("%+v", cfg)
	return &cfg, err
}
