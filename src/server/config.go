package server

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	User1      UserInfo `yaml:"user1"`
	User2      UserInfo `yaml:"user2"`
	User3      UserInfo `yaml:"user3"`
}

type HTTPServer struct {
	SSLcert     string        `yaml:"sslCert" env-default:""`
	SSLkey      string        `yaml:"sslKey"  env-default:""`
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type UserInfo struct {
	User_Email    string `yaml:"email" env-default:""`
	User_Name     string `yaml:"name" env-default:""`
	User_Password string `yaml:"password" env-default:""`
	User_Payload  string `yaml:"payload" env-default:""`
}

func ConfigLoad() *Config {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.yaml"
		log.Printf("Environment variable CONFIG_PATH is not set. Read default from %s.", configPath)
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Printf("Error opening config file: %v", err.Error())
		os.Exit(1)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Printf("Error reading config file: %v ", err.Error())
		os.Exit(1)
	}

	return &cfg
}
