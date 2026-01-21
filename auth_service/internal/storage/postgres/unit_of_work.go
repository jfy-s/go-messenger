package postgres

import (
	"context"
	"log/slog"
	"messenger-auth/internal/storage"

	"github.com/jackc/pgx/v5"
)

type UnitOfWork struct {
	tx       pgx.Tx
	userRepo UserRepository
	logger   *slog.Logger
}

func NewUnitOfWork(tx pgx.Tx, logger *slog.Logger) *UnitOfWork {
	return &UnitOfWork{tx: tx, userRepo: UserRepository{tx: tx, logger: logger}, logger: logger}
}

func (u *UnitOfWork) UserRepository() storage.UserRepository {
	return &u.userRepo
}

func (u *UnitOfWork) Commit() error {
	return u.tx.Commit(context.Background())
}

func (u *UnitOfWork) Rollback() error {
	return u.tx.Rollback(context.Background())
}
