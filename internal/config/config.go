package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	CoinGeckoAPIKey = iota
)

var conf *Config

type Config struct {
	CoinGeckoAPIKey string `yaml:"coin-gecko-api-key"`
}

func Init() error {
	ENV := os.Getenv("ENV")

	body, err := os.ReadFile(fmt.Sprintf("./configs/values_%s.yaml", strings.ToLower(ENV)))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &conf)
	return err
}

func Get(key int) interface{} {
	switch key {
	case CoinGeckoAPIKey:
		return conf.CoinGeckoAPIKey
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
