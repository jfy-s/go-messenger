package postgres

import (
	"context"
	"log/slog"
	"websocket_manager/internal/model"

	"github.com/jackc/pgx/v5"
)

type ChatRepository struct {
	tx     pgx.Tx
	logger *slog.Logger
}

func (repo *ChatRepository) GetAllUserChats(id uint64) ([]model.Chat, error) {
	rows, err := repo.tx.Query(context.Background(), "SELECT chat_id, name, creator_id FROM chat_users JOIN chats ON chats.id = chat_id WHERE user_id = $1", id)
	if err != nil {
		repo.logger.Error("failed to get chat ids", "error", err)
		return nil, err
	}

	chats := make([]model.Chat, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		var chat model.Chat
		err = rows.Scan(&chat.ID, &chat.Name, &chat.CreatorID)
		if err != nil {
			repo.logger.Error("failed to scan chat", "error", err)
			return nil, err
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (repo *ChatRepository) CreateChat(chat *model.Chat) error {
	err := repo.tx.QueryRow(context.Background(), "INSERT INTO chats (name, creator_id) VALUES ($1, $2) RETURNING id", chat.Name, chat.CreatorID).Scan(&chat.ID)
	if err != nil {
		repo.logger.Error("failed to create chat", "error", err)
	}

	return err
}

func (repo *ChatRepository) UpdateChat(chat *model.Chat) error {
	_, err := repo.tx.Exec(context.Background(), "UPDATE chats SET name = $1, updated_at = now() WHERE id = $2", chat.Name, chat.ID)
	if err != nil {
		repo.logger.Error("failed to update chat", "error", err)
	}

	return err
}

func (repo *ChatRepository) DeleteChat(id uint64) error {
	_, err := repo.tx.Exec(context.Background(), "DELETE FROM chats WHERE id = $1", id)
	if err != nil {
		repo.logger.Error("failed to delete chat", "error", err)
	}

	return err
}

func (repo *ChatRepository) AddUserToChat(chatUsers *model.ChatUsers) error {
	_, err := repo.tx.Exec(context.Background(), "INSERT INTO chat_users (chat_id, user_id) VALUES ($1, $2)", chatUsers.ChatID, chatUsers.UserID)
	if err != nil {
		repo.logger.Error("failed to add user to chat", "error", err)
	}

	return err
}

func (repo *ChatRepository) DeleteUserFromChat(chatUsers *model.ChatUsers) error {
	_, err := repo.tx.Exec(context.Background(), "DELETE FROM chat_users WHERE chat_id = $1 AND user_id = $2", chatUsers.ChatID, chatUsers.UserID)
	if err != nil {
		repo.logger.Error("failed to delete user from chat", "error", err)
	}

	return err
}

func (repo *ChatRepository) GetAllUsersIDInChat(id uint64) ([]uint64, error) {
	rows, err := repo.tx.Query(context.Background(), "SELECT user_id FROM chat_users WHERE chat_id = $1", id)
	if err != nil {
		repo.logger.Error("failed to get all users in chat", "error", err)
	}

	ids := make([]uint64, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			repo.logger.Error("failed to scan user id", "error", err)
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (repo *ChatRepository) GetOwnerID(id uint64) (uint64, error) {
	var ownerId uint64
	err := repo.tx.QueryRow(context.Background(), "SELECT creator_id FROM chats WHERE id = $1", id).Scan(&ownerId)
	if err != nil {
		repo.logger.Error("failed to get owner id", "error", err)
		return 0, err
	}
	return ownerId, nil
}

func (repo *ChatRepository) GetChatInfo(id uint64) (*model.Chat, []model.User, error) {
	chat := &model.Chat{ID: id}
	err := repo.tx.QueryRow(context.Background(), "SELECT name, creator_id FROM chats WHERE id = $1", id).Scan(&chat.Name, &chat.CreatorID)
	if err != nil {
		repo.logger.Error("failed to get chat info", "error", err)
		return nil, nil, err
	}

	rows, err := repo.tx.Query(context.Background(), "SELECT user_id, name FROM chat_users JOIN users ON user_id = users.id WHERE chat_id = $1", id)
	if err != nil {
		repo.logger.Error("failed to get chat info", "error", err)
		return nil, nil, err
	}
	users := make([]model.User, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			repo.logger.Error("failed to scan user id", "error", err)
			return nil, nil, err
		}
		users = append(users, user)
	}

	return chat, users, nil
}
