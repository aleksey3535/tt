package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	SecretKey string
	Env       string
	Db        DbConf
}

type DbConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Sslmode  string
}

func MustLoad() *Config {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH environment variable is not set")
	}
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading config file " + err.Error())
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic("Error unmarshalling config " + err.Error())
	}
	cfg.Db.Password = os.Getenv("DB_PASSWORD")
	cfg.SecretKey = os.Getenv("SECRET_KEY")
	return &cfg
}
