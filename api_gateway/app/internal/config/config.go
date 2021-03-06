package config

import (
	"github.com/WTC-SYSTEM/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	JWT     struct {
		Secret string `yaml:"secret" env:"JWT_SECRET" env-required:"true"`
	} `yaml:"jwt" env-required:"true"`
	Redis struct {
		Password string `yaml:"redis_password" env:"REDIS_PASSWORD" env-required:"true"`
		Addr     string `yaml:"redis_addr" env:"REDIS_ADDR" env-required:"true"`
	} `yaml:"redis" env-required:"true"`
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	UserService struct {
		URL string `yaml:"url" env-required:"true" env:"USER_SERVICE_URL"`
	} `yaml:"user_service" env-required:"true"`
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
