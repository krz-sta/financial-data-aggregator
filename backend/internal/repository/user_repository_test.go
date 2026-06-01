package repository

import (
	"context"
	"financial-data-aggregator-backend/internal/models"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestUserRepository(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	user1ID := uuid.New()
	user2ID := uuid.New()

	t.Run("Create User", func(t *testing.T) {
		user := &models.User{
			ID:           user1ID,
			Email:        "create@test.com",
			DisplayName:  "Create Test User",
			PasswordHash: "hashedpassword",
		}

		err := repo.Create(user)
		assert.NoError(t, err)
	})

	t.Run("Email Exists", func(t *testing.T) {
		exists, err := repo.EmailExists("create@test.com")
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.EmailExists("nobody@test.com")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Find By Email", func(t *testing.T) {
		user, err := repo.FindByEmail("create@test.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "Create Test User", user.DisplayName)
		assert.Equal(t, user1ID, user.ID)
	})

	t.Run("Find By ID", func(t *testing.T) {
		user2 := &models.User{
			ID:           user2ID,
			Email:        "findbyid@test.com",
			DisplayName:  "ID Test User",
			PasswordHash: "hashedpassword",
		}
		err := repo.Create(user2)
		assert.NoError(t, err)

		foundUser, err := repo.FindById(user2ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, "findbyid@test.com", foundUser.Email)
	})
}
