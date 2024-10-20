package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"local" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server"`
	DB         DataBase   `yaml:"db"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"jwt-auth-service"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DataBase struct {
	Username   string `yaml:"username" env-default:"postgres"`
	Host       string `yaml:"host" env-default:"db"`
	Port       string `yaml:"port" env-default:"5432"`
	DBName     string `yaml:"dbname" env-default:"my-db"`
	DBPassword string
	SSLMode    string `yaml:"sslmode" env-default:"disable"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is empty")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exits: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
	cfg.DB.DBPassword = os.Getenv("DB_PASSWORD")
	return &cfg
}
