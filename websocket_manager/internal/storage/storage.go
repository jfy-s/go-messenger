package storage

import "websocket_manager/internal/model"

type Storage interface {
	CreateUnitOfWork() (UnitOfWork, error)
	Close()
}

type UnitOfWork interface {
	ChatRepository() ChatRepository
	MessageRepository() MessageRepository
	Commit() error
	Rollback() error
}

type ChatRepository interface {
	CreateChat(chat *model.Chat) error
	UpdateChat(chat *model.Chat) error
	DeleteChat(id uint64) error
	AddUserToChat(chatUsers *model.ChatUsers) error
	DeleteUserFromChat(chatUsers *model.ChatUsers) error
	GetAllUsersIDInChat(id uint64) ([]uint64, error)
	GetOwnerID(id uint64) (uint64, error)
	GetAllUserChats(id uint64) ([]model.Chat, error)
	GetChatInfo(id uint64) (*model.Chat, []model.User, error)
}

type MessageRepository interface {
	AddMessage(msg *model.Message) error
	UpdateMessage(msg *model.Message) error
	DeleteMessage(id uint64) error
	GetAllMessagesInChat(chatID uint64) ([]model.Message, error)
	GetSenderID(id uint64) (uint64, error)
}
