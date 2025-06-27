package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func NewPostgresStorage(config string) (*Storage, error) {
	const op = "Storage.NewPostgresStorage"

	db, errOpen := sql.Open("postgres", config)
	if errOpen != nil {
		return nil, fmt.Errorf("%s : %w", op, errOpen)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Connection() *sql.DB {
	return s.db
}

func (s *Storage) SaveTrade(ctx context.Context, buyOrderId, sellOrderId, amount, price, timestamp int64) (int64, error) {
	const op = "Storage.Postgres.SaveTrade"

	var id int64
	err := s.db.QueryRowContext(ctx,
		"insert into Trades(buy_order_id, sell_order_id, amount, price, timestamp) values ($1, $2, $3, $4, $5) returning id",
		buyOrderId, sellOrderId, amount, price, timestamp).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	return id, nil
}
