package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name            string
	Version         string
	Port            int
	Env             string
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
	LogDir          string `mapstructure:"log_dir"`
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host          string
	Port          int
	Password      string
	DB            int
	PoolSize      int `mapstructure:"pool_size"`
	MinIdleConns  int `mapstructure:"min_idle_conns"`
	MaxRetries    int `mapstructure:"max_retries"`
	RetryInterval int `mapstructure:"retry_interval"`
}

type JWTConfig struct {
	Secret            string `mapstructure:"secret"`
	RefreshSecret     string `mapstructure:"refresh_secret"`
	Expiration        int    `mapstructure:"expiration"`         // Saat cinsinden
	RefreshExpiration int    `mapstructure:"refresh_expiration"` // Saat cinsinden
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config okuma hatası: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
