package config

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	minTokenTTL     = 10 * time.Minute
	maxTokenTTL     = 3 * time.Hour
	defaultTokenTTL = 15 * time.Minute
)

type Config struct {
	Port          int
	Jwt           Jwt
	AccessControl map[string][]string
	Dependencies  map[string]string
}
type Jwt struct {
	Secret []byte
	Ttl    time.Duration
}

func LoadConfig() *Config {
	cfg := &Config{
		Port: 50016,
		Jwt: Jwt{
			Ttl:    normalizeTime(viper.GetDuration("jwt.ttl")),
			Secret: []byte(viper.GetString("jwt.secret")),
		},
		AccessControl: viper.GetStringMapStringSlice("permissions"),
		Dependencies:  viper.GetStringMapString("dependencies.services"),
	}
	return cfg
}

func InitViper(filename string) {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./.conf")

	if err := viper.ReadInConfig(); err != nil {
		panic("error reading config file")
	}

	viper.OnConfigChange(func(e fsnotify.Event) {})
	viper.WatchConfig()
}

func normalizeTime(ttlTime time.Duration) time.Duration {
	var ttl = ttlTime

	if ttl <= 0 {
		return defaultTokenTTL
	}

	if ttl < minTokenTTL {
		ttl = minTokenTTL
	} else if ttl > maxTokenTTL {
		ttl = maxTokenTTL
	}
	return ttl

}
