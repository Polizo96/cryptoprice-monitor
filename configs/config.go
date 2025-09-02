package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type ExchangeConfig struct {
	Name      string `mapstructure:"name"`
	ApiUrl    string `mapstructure:"api_url"`
	ApiKey    string
	ApiSecret string
}

type SecretExchange struct {
	Name      string `mapstructure:"name"`
	ApiKey    string `mapstructure:"api_key"`
	ApiSecret string `mapstructure:"api_secret"`
}

type Secrets struct {
	Exchanges []SecretExchange `mapstructure:"exchanges"`
}

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Fetcher struct {
		Exchanges       []ExchangeConfig `mapstructure:"exchanges"`
		IntervalSeconds int              `mapstructure:"interval_seconds"`
	} `mapstructure:"fetcher"`
	Coins   []string `mapstructure:"coins"`
	Logging struct {
		File       string `mapstructure:"file"`       // path to log file, empty = stdout only
		EnableFile bool   `mapstructure:"enableFile"` // whether to enable file logging
	} `mapstructure:"logging"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	secretsViper := viper.New()
	secretsViper.SetConfigName("secrets")
	secretsViper.AddConfigPath("./configs")
	if err := secretsViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading secrets file: %w", err)
	}

	var secrets Secrets
	if err := secretsViper.Unmarshal(&secrets); err != nil {
		return nil, fmt.Errorf("unable to decode secrets: %w", err)
	}

	secretMap := make(map[string]SecretExchange)
	for _, se := range secrets.Exchanges {
		secretMap[se.Name] = se
	}

	for i, ex := range config.Fetcher.Exchanges {
		if sec, ok := secretMap[ex.Name]; ok {
			config.Fetcher.Exchanges[i].ApiKey = sec.ApiKey
			config.Fetcher.Exchanges[i].ApiSecret = sec.ApiSecret
		}
	}

	return &config, nil
}
