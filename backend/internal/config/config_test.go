package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("DB_HOST", "test_host")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_NAME", "test_name")
	os.Setenv("DB_PORT", "test_port")
	os.Setenv("ROUTER_HOST", "test_router_host")
	os.Setenv("ROUTER_PORT", "test_router_port")
	os.Setenv("JWT_SECRET", "test_jwt_secret")

	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")
	defer os.Unsetenv("DB_NAME")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("ROUTER_HOST")
	defer os.Unsetenv("ROUTER_PORT")
	defer os.Unsetenv("JWT_SECRET")

	cfg := LoadConfig()
	dsn := cfg.DB.GetDsn()
	rAddr := cfg.Router.GetRouterConfig()
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
}
