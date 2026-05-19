package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB     *DBconfig
	Router *RouterConfig
	JWTKey string
}

type DBconfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type RouterConfig struct {
	Host string
	Port string
}

func (db *DBconfig) GetDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		db.Host, db.User, db.Password, db.Name, db.Port)
}

func (r *RouterConfig) GetRouterConfig() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func LoadConfig() *Config {
	dbCfg := &DBconfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}

	routerCfg := &RouterConfig{
		Host: os.Getenv("ROUTER_HOST"),
		Port: os.Getenv("ROUTER_PORT"),
	}

	return &Config{
		DB:     dbCfg,
		Router: routerCfg,
		JWTKey: os.Getenv("JWT_SECRET"),
	}
}
