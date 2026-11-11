package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrMissingEnvironmentVariable = fmt.Errorf("missing required environment variable")
	ErrInvalidEnvironmentVariable = fmt.Errorf("invalid environment variable value")
)

type Config struct {
	Env   string      // Environment，区分开发环境、测试环境和生产环境
	MySQL MySQLConfig // MySQL配置
	HTTP  HTTPConfig  // HTTP配置
}

type MySQLConfig struct {
	DSN             string        // Data Source Name，MySQL连接字符串
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
}

type HTTPConfig struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func Load() (Config, error) {
	var cfg Config

	env, exist := os.LookupEnv("APP_ENV")
	if !exist {
		cfg.Env = "development"
	} else {
		cfg.Env = env
	}

	dsn, exist := os.LookupEnv("MYSQL_DSN")
	if !exist {
		return cfg, fmt.Errorf("MYSQL_DSN is required: %w", ErrMissingEnvironmentVariable)
	}
	cfg.MySQL.DSN = dsn

	maxOpenConns, exist := os.LookupEnv("MYSQL_MAX_OPEN_CONNS")
	if !exist {
		cfg.MySQL.MaxOpenConns = 10
	} else {
		val, err := strconv.Atoi(maxOpenConns)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse MYSQL_MAX_OPEN_CONNS: %w", ErrInvalidEnvironmentVariable)
		}
		cfg.MySQL.MaxOpenConns = val
	}

	maxIdleConns, exist := os.LookupEnv("MYSQL_MAX_IDLE_CONNS")
	if !exist {
		cfg.MySQL.MaxIdleConns = 5
	} else {
		val, err := strconv.Atoi(maxIdleConns)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse MYSQL_MAX_IDLE_CONNS: %w", ErrInvalidEnvironmentVariable)
		}
		cfg.MySQL.MaxIdleConns = val
	}

	connMaxLifetime, exist := os.LookupEnv("MYSQL_CONN_MAX_LIFETIME")
	if !exist {
		cfg.MySQL.ConnMaxLifetime = 30 * time.Second
	} else {
		duration, err := time.ParseDuration(connMaxLifetime)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse MYSQL_CONN_MAX_LIFETIME:  %w", ErrInvalidEnvironmentVariable)
		}
		cfg.MySQL.ConnMaxLifetime = duration
	}

	addr, exist := os.LookupEnv("HTTP_ADDR")
	if !exist {
		cfg.HTTP.Addr = ":8080"
	} else {
		cfg.HTTP.Addr = addr
	}

	readHeaderTimeout, exist := os.LookupEnv("HTTP_READ_HEADER_TIMEOUT")
	if !exist {
		cfg.HTTP.ReadHeaderTimeout = 5 * time.Second
	} else {
		duration, err := time.ParseDuration(readHeaderTimeout)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse HTTP_READ_HEADER_TIMEOUT: %w", ErrInvalidEnvironmentVariable)
		}
		cfg.HTTP.ReadHeaderTimeout = duration
	}

	readTimeout, exist := os.LookupEnv("HTTP_READ_TIMEOUT")
	if !exist {
		cfg.HTTP.ReadTimeout = 15 * time.Second
	} else {
		duration, err := time.ParseDuration(readTimeout)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse HTTP_READ_TIMEOUT: %w", ErrInvalidEnvironmentVariable)
		}
		cfg.HTTP.ReadTimeout = duration
	}

	writeTimeout, exist := os.LookupEnv("HTTP_WRITE_TIMEOUT")
	if !exist {
		cfg.HTTP.WriteTimeout = 15 * time.Second
	} else {
		duration, err := time.ParseDuration(writeTimeout)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse HTTP_WRITE_TIMEOUT: %w", ErrInvalidEnvironmentVariable)
		}
		cfg.HTTP.WriteTimeout = duration
	}

	idleTimeout, exist := os.LookupEnv("HTTP_IDLE_TIMEOUT")
	if !exist {
		cfg.HTTP.IdleTimeout = 60 * time.Second
	} else {
		duration, err := time.ParseDuration(idleTimeout)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse HTTP_IDLE_TIMEOUT: %w", ErrInvalidEnvironmentVariable)
		}
		cfg.HTTP.IdleTimeout = duration
	}

	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.Env) == "" {
		return fmt.Errorf("APP_ENV is required: %w", ErrInvalidEnvironmentVariable)
	}

	if cfg.Env != "development" && cfg.Env != "test" && cfg.Env != "production" {
		return fmt.Errorf("invalid APP_ENV %q: %w", cfg.Env, ErrInvalidEnvironmentVariable)
	}
	if cfg.MySQL.MaxOpenConns <= 0 {
		return fmt.Errorf("MYSQL_MAX_OPEN_CONNS must be greater than 0: %w", ErrInvalidEnvironmentVariable)
	}
	if cfg.MySQL.MaxIdleConns < 0 {
		return fmt.Errorf("MYSQL_MAX_IDLE_CONNS must be greater than or equal to 0: %w", ErrInvalidEnvironmentVariable)
	}
	if cfg.MySQL.MaxOpenConns < cfg.MySQL.MaxIdleConns {
		return fmt.Errorf("MYSQL_MAX_OPEN_CONNS must be greater than or equal to MYSQL_MAX_IDLE_CONNS: %w", ErrInvalidEnvironmentVariable)
	}

	if cfg.MySQL.ConnMaxLifetime <= 0 {
		return fmt.Errorf("MYSQL_CONN_MAX_LIFETIME must be greater than 0: %w", ErrInvalidEnvironmentVariable)
	}

	if strings.TrimSpace(cfg.HTTP.Addr) == "" {
		return fmt.Errorf("HTTP_ADDR is required: %w", ErrInvalidEnvironmentVariable)
	}
	if cfg.HTTP.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("HTTP_READ_HEADER_TIMEOUT must be greater than 0: %w", ErrInvalidEnvironmentVariable)
	}
	if cfg.HTTP.ReadTimeout <= 0 {
		return fmt.Errorf("HTTP_READ_TIMEOUT must be greater than 0: %w", ErrInvalidEnvironmentVariable)
	}
	if cfg.HTTP.WriteTimeout <= 0 {
		return fmt.Errorf("HTTP_WRITE_TIMEOUT must be greater than 0: %w", ErrInvalidEnvironmentVariable)
	}
	if cfg.HTTP.IdleTimeout <= 0 {
		return fmt.Errorf("HTTP_IDLE_TIMEOUT must be greater than 0: %w", ErrInvalidEnvironmentVariable)
	}
	return nil
}
