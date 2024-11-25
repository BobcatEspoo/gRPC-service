package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
	"time"
)

type RateService struct {
	Db *sql.DB
	grpc_health_v1.UnimplementedHealthServer
	UnimplementedRatesServiceServer
}

func NewRateService(dbService *sql.DB) *RateService {
	return &RateService{Db: dbService}
}

func (ds *RateService) GetRates(_ context.Context, req *GetRatesRequest) (res *GetRatesResponse, err error) {
	resp, err := http.Get("https://garantex.org/api/v2/depth?market=" + req.Market)
	if err != nil {
		return nil, fmt.Errorf("problem in connection to garantex.org:%w", err)
	}
	defer resp.Body.Close()
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("problem in decoding body from garantex.org:%w", err)
	}
	date := time.Unix(response.Timestamp, 0).Format("2006-01-02 15:04:05")
	err = ds.AddToDatabase(req.Market, date, response.Asks[0])
	if err != nil {
		return nil, fmt.Errorf("problem in adding to database:%w", err)
	}
	res = &GetRatesResponse{
		Time: date,
		Asks: response.Asks[0],
	}
	return res, nil
}

type Response struct {
	Timestamp int64   `json:"timestamp"`
	Asks      []*Asks `json:"asks"`
}

func (r *RateService) AddToDatabase(request, date string, asks *Asks) error {
	query := "INSERT INTO answers (request, time, price, volume, amount, factor, type) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.Db.Exec(query, request, date, asks.Price, asks.Volume, asks.Amount, asks.Factor, asks.Type)
	if err != nil {
		return fmt.Errorf("error inserting data into database: %w", err)
	}
	return nil
}

func (ds *RateService) mustEmbedUnimplementedRatesServiceServer() {}
