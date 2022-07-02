package config

import (
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
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
	AwsCfg  AwsConfig     `yaml:"aws"`
}

type StorageConfig struct {
	Host        string `yaml:"host" env:"DB_HOST" env-required:"true"`
	Port        string `yaml:"port" env:"DB_PORT" env-required:"true"`
	Database    string `yaml:"database" env:"DB_DATABASE" env-required:"true"`
	Username    string `yaml:"username" env:"DB_USERNAME" env-required:"true"`
	Password    string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	MaxAttempts int8   `yaml:"maxAttempts" env:"DB_MAX_ATTEMPTS" env-default:"5"`
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
