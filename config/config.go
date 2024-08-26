package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	MySQLConfig `yaml:"mysql"`
	SecretKey   string `yaml:"secret_key"`
}

var once sync.Once

var CFG = GetConfig()

func GetConfig() *Config {
	cfg := &Config{}
	once.Do(func() {
		log.Println("Reading app configuration...")
		if err := cleanenv.ReadConfig("/notes/config/config.yml", cfg); err != nil {
			log.Fatalln(err)
		}
		log.Println("Successfully loaded config")
	})
	return cfg
}
