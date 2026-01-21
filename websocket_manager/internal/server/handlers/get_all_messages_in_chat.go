package handlers

import (
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type GetAllMessagesInChatRequest struct {
	UserID uint64 `validate:"required, min=1"`
	ChatID uint64 `validate:"required"`
}

func HandleGetAllMessagesInChat(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	req := GetAllMessagesInChatRequest{UserID: msgPacketRequest.From, ChatID: msgPacketRequest.To}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllMessagesInChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllMessagesInChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	defer uow.Rollback()
	msgRepo := uow.MessageRepository()
	msgs := make([]model.Message, 0)
	msgs, err = msgRepo.GetAllMessagesInChat(req.ChatID)
	if err != nil {
		logger.Error("failed to get messages", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllMessagesInChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllMessagesInChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	logger.Info("messages received", "count", len(msgs), "chat_id", req.ChatID, "user_id", req.UserID)
	return &model.MessagePacketRequest{MsgType: model.GetAllMessagesInChat, From: 0, To: msgPacketRequest.From, Data: "Success"}
}
