package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

const SERVICENAME = "multiplexer"

type (
	HttpConfig struct {
		Port             int           `json:"port",mapstructure:"http_port"`
		Timeout          time.Duration `json:"timeout",mapstructure:"http_timeout"`
		ConnectionsLimit int           `json:"timeout",mapstructure:"http_connectionsLimit"`
	}
	ClientConfig struct {
		PerRequestTimeout time.Duration `json:"timeout",mapstructure:"client_perRequestTimeout"`
		MaxWorkers        uint8         `json:"max_workers",mapstructure:"client_maxWorkers"`
		MaxUrls           int           `json:"max_workers",mapstructure:"maxUrls"`
	}

	RootConfig struct {
		Http   HttpConfig
		Client ClientConfig
	}
)

func MustConfigure() *RootConfig {
	viper.SetEnvPrefix(SERVICENAME)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.timeout", "10s")
	viper.SetDefault("http.connectionsLimit", 100)
	viper.SetDefault("client.perRequestTimeout", "1s")
	viper.SetDefault("client.maxWorkers", 4)
	viper.SetDefault("client.maxUrls", 20)

	cfg := &RootConfig{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return cfg
	}
	return cfg
}
