package postgres

import (
	"context"
	"log/slog"
	"messenger-auth/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewStorage(databaseUrl string, logger *slog.Logger) (*Storage, error) {
	db, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, err

	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &Storage{db: db, logger: logger}, err
}

func (s *Storage) CreateUnitOfWork() (storage.UnitOfWork, error) {
	tx, err := s.db.Begin(context.Background())

	if err != nil {
		s.logger.Error("failed start transaction", "error", err)
		return nil, err
	}

	return NewUnitOfWork(tx, s.logger), nil
}

func (s *Storage) Close() {
	s.db.Close()
}
