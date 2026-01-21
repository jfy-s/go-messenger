package postgres

import (
	"context"
	"log/slog"
	"websocket_manager/internal/storage"

	"github.com/jackc/pgx/v5"
)

type UnitOfWork struct {
	tx          pgx.Tx
	chatRepo    ChatRepository
	messageRepo MessageRepository
	logger      *slog.Logger
}

func NewUnitOfWork(tx pgx.Tx, logger *slog.Logger) *UnitOfWork {
	return &UnitOfWork{tx: tx, chatRepo: ChatRepository{tx: tx, logger: logger}, messageRepo: MessageRepository{tx: tx, logger: logger}, logger: logger}
}

func (u *UnitOfWork) ChatRepository() storage.ChatRepository {
	return &u.chatRepo
}

func (u *UnitOfWork) MessageRepository() storage.MessageRepository {
	return &u.messageRepo
}

func (u *UnitOfWork) Commit() error {
	return u.tx.Commit(context.Background())
}

func (u *UnitOfWork) Rollback() error {
	return u.tx.Rollback(context.Background())
}
