package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	CoinGeckoAPIKey = iota
	CacheTTL
	RedisHost
	RedisPassword
)

var conf *Config

type Config struct {
	CoinGeckoAPIKey string        `yaml:"coin-gecko-api-key"`
	CacheTTL        time.Duration `yaml:"cache-ttl"`
	RedisHost       string        `yaml:"redis-host"`
	RedisPassword   string        `yaml:"redis-password"`
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
	case CacheTTL:
		return conf.CacheTTL
	case RedisHost:
		return conf.RedisHost
	case RedisPassword:
		return conf.RedisPassword
	default:
		panic(ErrConfigNotFoundByKey(key))
	}
}
