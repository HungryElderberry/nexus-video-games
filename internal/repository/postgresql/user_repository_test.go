package postgresql_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"nexus-video-games/internal/entity"
	"nexus-video-games/internal/repository/postgresql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := postgresql.NewUserRepository(gormDB)
	userID := uuid.New()
	email := "gamer@nexus.com"

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "is_verified", "created_at", "updated_at"}).
		AddRow(userID, email, "hashed_pw", true, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE email = \$1 .*LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnRows(rows)

	result, err := repo.FindByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	if result != nil {
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, email, result.Email)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := postgresql.NewUserRepository(gormDB)
	email := "unknown@nexus.com"

	mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE email = \$1 .*LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.FindByEmail(context.Background(), email)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Create_Success(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := postgresql.NewUserRepository(gormDB)
	userID := uuid.New()
	newUser := &entity.User{
		ID:           userID,
		Email:        "newgamer@nexus.com",
		PasswordHash: "secret123",
		IsVerified:   false,
	}

	returningRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(userID, time.Now(), time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(newUser.Email, newUser.PasswordHash, newUser.IsVerified, newUser.ID).
		WillReturnRows(returningRows)
	mock.ExpectCommit()

	err := repo.Create(context.Background(), newUser)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
