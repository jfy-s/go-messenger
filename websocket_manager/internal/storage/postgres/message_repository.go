package postgres

import (
	"context"
	"log/slog"
	"websocket_manager/internal/model"

	"github.com/jackc/pgx/v5"
)

type MessageRepository struct {
	tx     pgx.Tx
	logger *slog.Logger
}

func (repo *MessageRepository) AddMessage(msg *model.Message) error {
	err := repo.tx.QueryRow(context.Background(), "INSERT INTO messages (chat_id, user_id, message) VALUES ($1, $2, $3) RETURNING id", msg.ChatID, msg.UserID, msg.Message).Scan(&msg.ID)
	if err != nil {
		repo.logger.Error("failed to send message", "error", err)
	}

	return err
}

func (repo *MessageRepository) UpdateMessage(msg *model.Message) error {
	_, err := repo.tx.Exec(context.Background(), "UPDATE messages SET message = $1, updated_at = now() WHERE id = $2", msg.Message, msg.ID)
	if err != nil {
		repo.logger.Error("failed to update message", "error", err)
	}

	return err
}

func (repo *MessageRepository) DeleteMessage(id uint64) error {
	_, err := repo.tx.Exec(context.Background(), "DELETE FROM messages WHERE id = $1", id)
	if err != nil {
		repo.logger.Error("failed to delete message", "error", err)
	}

	return err
}

func (repo *MessageRepository) GetAllMessagesInChat(chatID uint64) ([]model.Message, error) {
	rows, err := repo.tx.Query(context.Background(), "SELECT id, user_id, message, created_at, updated_at FROM messages WHERE chat_id = $1", chatID)
	if err != nil {
		repo.logger.Error("failed to get all messages in chat", "error", err)
	}

	msgs := make([]model.Message, rows.CommandTag().RowsAffected())
	for rows.Next() {
		msg := model.Message{ChatID: chatID}
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Message, &msg.CreatedAt, &msg.UpdatedAt); err != nil {
			repo.logger.Error("failed to scan message", "error", err)
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (repo *MessageRepository) GetSenderID(id uint64) (uint64, error) {
	var senderID uint64
	err := repo.tx.QueryRow(context.Background(), "SELECT user_id FROM messages WHERE id = $1", id).Scan(&senderID)
	if err != nil {
		repo.logger.Error("failed to get sender id", "error", err)
	}
	return senderID, nil
}
