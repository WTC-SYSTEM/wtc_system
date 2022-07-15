package config

import (
	"github.com/WTC-SYSTEM/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	AwsCfg AwsConfig `yaml:"aws"`
}

type AwsConfig struct {
	AccessKeyID     string `yaml:"accessKeyID" env:"AWS_ACCESS_KEY_ID" env-required:"true"`
	SecretAccessKey string `yaml:"secretAccessKey" env:"AWS_SECRET_ACCESS_KEY" env-required:"true"`
	Region          string `yaml:"region" env:"AWS_REGION" env-required:"true"`
	Bucket          string `yaml:"bucket" env:"AWS_BUCKET" env-required:"true"`
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
