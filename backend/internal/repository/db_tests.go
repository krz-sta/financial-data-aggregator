package repository

import (
	"context"
	"financial-data-aggregator-backend/internal/models"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	postgresContainer, err := tcpostgres.Run(ctx,
		"postgres:15-alpine",
		tcpostgres.WithDatabase("test_db"),
		tcpostgres.WithUsername("test_user"),
		tcpostgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(time.Second*10),
		),
	)
	if err != nil {
		t.Fatalf("failed to start posgres container: %v", err)
	}

	testDsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get dsn string: %v", err)
	}

	db, err := gorm.Open(gormpostgres.Open(testDsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.PortfolioItem{})
	if err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}

	cleanup := func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate the test databse: %v", err)
		}
	}

	return db, cleanup
}
