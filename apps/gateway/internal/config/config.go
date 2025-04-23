package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env     string        `yaml:"env" env-required:"true"`
	Port    string        `yaml:"port" env-default:":8080"`
	Clients ClientsConfig `yaml:"clients"`
}

type Auth struct {
	Address   string `yaml:"address"`
	Port      string `yaml:"port" default:"50051"`
	AppSecret string `yaml:"appSecret" env-required:"true" env:"APP_SECRET"`
}

type Products_customer struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port" default:"50052"`
}

type Products_admin struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port" default:"50053"`
}

type ClientsConfig struct {
	Auth              Auth              `yaml:"auth"`
	Products_customer Products_customer `yaml:"products-customer"`
	Products_admin    Products_admin    `yaml:"products-admin"`
}

func MustLoadConfig() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
