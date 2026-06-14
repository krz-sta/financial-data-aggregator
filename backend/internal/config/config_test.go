package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	t.Setenv("DB_HOST", "test_host")
	t.Setenv("DB_USER", "test_user")
	t.Setenv("DB_PASSWORD", "test_password")
	t.Setenv("DB_NAME", "test_name")
	t.Setenv("DB_PORT", "test_port")
	t.Setenv("ROUTER_HOST", "test_router_host")
	t.Setenv("ROUTER_PORT", "test_router_port")
	t.Setenv("JWT_SECRET", "test_jwt_secret")
	t.Setenv("REDIS_HOST", "test_redis_host")
	t.Setenv("REDIS_PORT", "test_redis_port")
	t.Setenv("REDIS_PASSWORD", "test_redis_passoword")

	// defer os.Unsetenv("DB_HOST")
	// defer os.Unsetenv("DB_USER")
	// defer os.Unsetenv("DB_PASSWORD")
	// defer os.Unsetenv("DB_NAME")
	// defer os.Unsetenv("DB_PORT")
	// defer os.Unsetenv("ROUTER_HOST")
	// defer os.Unsetenv("ROUTER_PORT")
	// defer os.Unsetenv("JWT_SECRET")
	// defer os.Unsetenv("REDIS_HOST")
	// defer os.Unsetenv("REDIS_PORT")
	// defer os.Unsetenv("REDIS_PASSWORD")

	cfg := LoadConfig()
	dsn := cfg.DB.GetDsn()
	rAddr := cfg.Router.GetRouterConfig()
	redAddr := cfg.Redis.GetRedisConfig()

	tJwt := cfg.JWTKey

	expectedDsn := "host=test_host user=test_user password=test_password dbname=test_name port=test_port sslmode=disable"
	if dsn != expectedDsn {
		t.Errorf("Expected %s, got %s", expectedDsn, dsn)
	}

	expectedRAddr := "test_router_host:test_router_port"
	if rAddr != expectedRAddr {
		t.Errorf("Expected %s, got %s", expectedRAddr, rAddr)
	}

	expectedJWTKey := "test_jwt_secret"
	if tJwt != expectedJWTKey {
		t.Errorf("Expected %s, got %s", expectedJWTKey, tJwt)
	}

	expectedRedAddr := "test_redis_host:test_redis_port"
	if redAddr != expectedRedAddr {
		t.Errorf("Expected %s, got %s", expectedRedAddr, redAddr)
	}
}
