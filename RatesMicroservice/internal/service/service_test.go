package service

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestRateService_AddToDatabase(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	rateService := NewRateService(db)

	request := "btc_usdt"
	date := "2024-11-23 15:04:05"
	asks := &Asks{
		Price:  "20000",
		Volume: "0.1",
		Amount: "2000",
		Factor: "10",
		Type:   "sell",
	}
	t.Run("success case", func(t *testing.T) {
		mock.ExpectExec(`^INSERT INTO answers \(request, time, price, volume, amount, factor, type\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)$`).
			WithArgs(request, date, asks.Price, asks.Volume, asks.Amount, asks.Factor, asks.Type).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := rateService.AddToDatabase(request, date, asks)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %v", err)
		}
	})

	t.Run("error case", func(t *testing.T) {
		mock.ExpectExec(`^INSERT INTO answers \(request, time, price, volume, amount, factor, type\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)$`).
			WithArgs(request, date, asks.Price, asks.Volume, asks.Amount, asks.Factor, asks.Type).
			WillReturnError(fmt.Errorf("mocked database error"))

		err := rateService.AddToDatabase(request, date, asks)
		if err == nil {
			t.Errorf("expected an error, but got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %v", err)
		}
	})

}
