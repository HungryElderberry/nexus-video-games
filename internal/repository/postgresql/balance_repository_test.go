package postgresql_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"nexus-video-games/internal/repository/postgresql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := gormPostgres.New(gormPostgres.Config{
		Conn:       dbMock,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		dbMock.Close()
	}

	return gormDB, mock, cleanup
}

func TestBalanceRepository_FindByUserID_Success(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := postgresql.NewBalanceRepository(gormDB)

	userID := uuid.New()
	expectedAmount := int64(500000)
	expectedCurrency := "IDR"

	rows := sqlmock.NewRows([]string{"user_id", "balance_amount", "currency", "updated_at"}).
		AddRow(userID, expectedAmount, expectedCurrency, time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM "user_balances" WHERE user_id = \$1 .*LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	result, err := repo.FindByUserID(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	if result != nil {
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, expectedAmount, result.BalanceAmount)
		assert.Equal(t, expectedCurrency, result.Currency)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBalanceRepository_FindByUserID_NotFound(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := postgresql.NewBalanceRepository(gormDB)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT (.+) FROM "user_balances" WHERE user_id = \$1 .*LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.FindByUserID(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBalanceRepository_AddBalance_Success(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := postgresql.NewBalanceRepository(gormDB)
	userID := uuid.New()
	initialAmount := int64(200000)
	addedAmount := int64(250000)

	rows := sqlmock.NewRows([]string{"user_id", "balance_amount", "currency", "updated_at"}).
		AddRow(userID, initialAmount, "IDR", time.Now())

	mock.ExpectQuery(`SELECT (.+) FROM "user_balances" WHERE user_id = \$1 .*LIMIT \$2 FOR UPDATE`).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "user_balances" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.AddBalance(context.Background(), gormDB, userID, addedAmount)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
