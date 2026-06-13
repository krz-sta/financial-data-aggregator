package repository

import (
	"financial-data-aggregator-backend/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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
