package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB     *DBconfig
	Router *RouterConfig
	Redis  *RedisConfig
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

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func (db *DBconfig) GetDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		db.Host, db.User, db.Password, db.Name, db.Port)
}

func (r *RouterConfig) GetRouterConfig() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (red *RedisConfig) GetRedisConfig() string {
	return fmt.Sprintf("%s:%s", red.Host, red.Port)
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

	redisCfg := &RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	return &Config{
		DB:     dbCfg,
		Router: routerCfg,
		Redis:  redisCfg,
		JWTKey: os.Getenv("JWT_SECRET"),
	}
}
