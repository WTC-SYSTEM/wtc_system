package config

import (
	"github.com/hawkkiller/wtc_system/user_service/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	Storage StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	Host        string `yaml:"host" env:"DB_HOST"`
	Port        string `yaml:"port" env:"DB_PORT"`
	Database    string `yaml:"database" env:"DB_NAME"`
	Username    string `yaml:"username" env:"DB_USERNAME"`
	Password    string `yaml:"password" env:"DB_PASSWORD"`
	MaxAttempts int8   `yaml:"maxAttempts" env-default:"5"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
