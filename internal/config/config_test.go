package config

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// 设置环境变量
	os.Setenv("APP_ENV", "test")
	os.Setenv("MYSQL_DSN", "user:password@tcp(localhost:3306)/testdb")
	os.Setenv("MYSQL_MAX_OPEN_CONNS", "20")
	os.Setenv("MYSQL_MAX_IDLE_CONNS", "10")
	os.Setenv("MYSQL_CONN_MAX_LIFETIME", "30s")
	os.Setenv("HTTP_ADDR", ":8080")
	os.Setenv("HTTP_READ_HEADER_TIMEOUT", "5s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "30s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Env != "test" {
		t.Errorf("expected Env to be 'test', got '%s'", cfg.Env)
	}
	if cfg.MySQL.DSN != "user:password@tcp(localhost:3306)/testdb" {
		t.Errorf("expected MySQL DSN to be 'user:password@tcp(localhost:3306)/testdb', got '%s'", cfg.MySQL.DSN)
	}
	if cfg.MySQL.MaxOpenConns != 20 {
		t.Errorf("expected MySQL MaxOpenConns to be 20, got %d", cfg.MySQL.MaxOpenConns)
	}
	if cfg.MySQL.MaxIdleConns != 10 {
		t.Errorf("expected MySQL MaxIdleConns to be 10, got %d", cfg.MySQL.MaxIdleConns)
	}
	if cfg.MySQL.ConnMaxLifetime != 30*time.Second {
		t.Errorf("expected MySQL ConnMaxLifetime to be 30s, got %v", cfg.MySQL.ConnMaxLifetime)
	}
	if cfg.HTTP.Addr != ":8080" {
		t.Errorf("expected HTTP Addr to be ':8080', got '%s'", cfg.HTTP.Addr)
	}
	if cfg.HTTP.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("expected HTTP ReadHeaderTimeout to be 5s, got %v", cfg.HTTP.ReadHeaderTimeout)
	}
	if cfg.HTTP.ReadTimeout != 10*time.Second {
		t.Errorf("expected HTTP ReadTimeout to be 10s, got %v", cfg.HTTP.ReadTimeout)
	}
	if cfg.HTTP.WriteTimeout != 10*time.Second {
		t.Errorf("expected HTTP WriteTimeout to be 10s, got %v", cfg.HTTP.WriteTimeout)
	}
	if cfg.HTTP.IdleTimeout != 30*time.Second {
		t.Errorf("expected HTTP IdleTimeout to be 30s, got %v", cfg.HTTP.IdleTimeout)
	}
}

func TestLoadMissingEnvironment(t *testing.T) {
	// 清除环境变量
	os.Unsetenv("APP_ENV")
	os.Unsetenv("MYSQL_DSN")
	os.Setenv("MYSQL_MAX_OPEN_CONNS", "20")
	os.Setenv("MYSQL_MAX_IDLE_CONNS", "10")
	os.Setenv("MYSQL_CONN_MAX_LIFETIME", "30s")
	os.Setenv("HTTP_ADDR", ":8080")
	os.Setenv("HTTP_READ_HEADER_TIMEOUT", "5s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "30s")
	_, err := Load()
	if !errors.Is(err, ErrMissingEnvironmentVariable) {
		t.Fatalf("expected ErrMissingEnvironmentVariable, got %v", err)
	}
}

func TestLoadInvalidEnvironment(t *testing.T) {
	// 设置无效的环境变量
	os.Setenv("APP_ENV", "test")
	os.Setenv("MYSQL_DSN", "user:password@tcp(localhost:3306)/testdb")
	os.Setenv("MYSQL_MAX_OPEN_CONNS", "abc") // 无效的整数
	os.Setenv("MYSQL_MAX_IDLE_CONNS", "10")
	os.Setenv("MYSQL_CONN_MAX_LIFETIME", "30s")
	os.Setenv("HTTP_ADDR", ":8080")
	os.Setenv("HTTP_READ_HEADER_TIMEOUT", "5s")
	os.Setenv("HTTP_READ_TIMEOUT", "10s")
	os.Setenv("HTTP_WRITE_TIMEOUT", "10s")
	os.Setenv("HTTP_IDLE_TIMEOUT", "30s")
	_, err := Load()
	if !errors.Is(err, ErrInvalidEnvironmentVariable) {
		t.Fatalf("expected ErrInvalidEnvironmentVariable, got %v", err)
	}
}

func TestValidateInvalidEnvironment(t *testing.T) {
	err := (Config{Env: "staging"}).Validate()
	if !errors.Is(err, ErrInvalidEnvironmentVariable) {
		t.Fatalf("expected ErrInvalidEnvironmentVariable, got %v", err)
	}
}
