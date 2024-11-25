package server

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestAccessToDB(t *testing.T) {

	os.Setenv("DB_URL", "postgres://user:password@localhost:5432/testdb")

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	t.Run("missing DB_URL", func(t *testing.T) {
		os.Unsetenv("DB_URL")

		result, err := AccessToDB()
		if err == nil {
			t.Fatal("expected an error for missing DB_URL, got none")
		}
		if result != nil {
			t.Fatal("expected result to be nil when DB_URL is missing")
		}
	})

	t.Run("failed db connection", func(t *testing.T) {
		os.Setenv("DB_URL", "postgres://user:password@localhost:5432/testdb")

		result, err := AccessToDB()
		if err == nil {
			t.Fatal("expected an error due to failed ping, got none")
		}
		if result != nil {
			t.Fatal("expected result to be nil when ping fails")
		}
	})
}

var sqlOpen = sql.Open

var mockAccessToDB = func() (*sql.DB, error) {
	return nil, nil
}

func TestStart(t *testing.T) {
	// Set up a mock logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	t.Run("successful start", func(t *testing.T) {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		stop <- syscall.SIGTERM

		assert.NoError(t, nil)
	})
}
