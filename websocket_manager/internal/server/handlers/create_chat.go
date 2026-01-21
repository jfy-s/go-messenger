package handlers

import (
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type CreateChatRequest struct {
	CreatorID uint64 `validate:"required, min=1"`
	Name      string `validate:"required,min=1,max=64"`
}

func HandleCreateChat(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	req := CreateChatRequest{CreatorID: msgPacketRequest.From, Name: msgPacketRequest.Data}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.CreateChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.CreateChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	defer uow.Rollback()
	chatRepo := uow.ChatRepository()
	chat := &model.Chat{CreatorID: req.CreatorID, Name: req.Name}
	err = chatRepo.CreateChat(chat)
	if err != nil {
		logger.Error("failed to create chat", "error", err)
		return &model.MessagePacketRequest{MsgType: model.CreateChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	chatUser := &model.ChatUsers{ChatID: chat.ID, UserID: req.CreatorID}
	err = chatRepo.AddUserToChat(chatUser)
	if err != nil {
		logger.Error("failed to add user to chat", "error", err)
		return &model.MessagePacketRequest{MsgType: model.CreateChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.CreateChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	logger.Info("chat created", "chat_id", chat.ID)
	return &model.MessagePacketRequest{MsgType: model.CreateChat, From: 0, To: msgPacketRequest.From, Data: "Success"}
}
