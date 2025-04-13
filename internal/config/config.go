package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local" `
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	DBConfig    `yaml:"db"`
}

type DBConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Username string `yaml:"user" env:"USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {

	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	config.DBConfig.Username = os.Getenv("DB_USER")
	config.DBConfig.Password = os.Getenv("DB_PASSWORD")
	config.DBConfig.Host = os.Getenv("DB_HOST")
	config.DBConfig.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	config.DBConfig.Name = os.Getenv("DB_NAME")

	return &config
}
