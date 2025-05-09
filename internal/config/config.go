package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
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
	configPath, err := fetchConfigPath()
	if err != nil {
		panic(err)
	} else if configPath == "" {
		log.Fatal("CONFIG_PATH is empty")
	} else if _, err = os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists")
	}

	var cfg Config
	if err = cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	err = godotenv.Load()
	if err != nil {
		panic("could not load .env")
	}
	cfg.DB.DBPassword = os.Getenv("DB_PASSWORD")

	return &cfg
}

func fetchConfigPath() (string, error) {
	var configPath string
	flag.StringVar(&configPath, "config_path", "", "path to config file")
	flag.Parse()
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath == "" {
		return "", fmt.Errorf("config path is empty")
	}

	return configPath, nil
}
