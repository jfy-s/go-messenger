package handlers

import (
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type SendMessageRequest struct {
	SenderID uint64 `validate:"required, min=1"`
	ChatID   uint64 `validate:"required"`
	Message  string `validate:"required"`
}

func HandleSendMessage(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	req := SendMessageRequest{SenderID: msgPacketRequest.From, ChatID: msgPacketRequest.To, Message: msgPacketRequest.Data}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	defer uow.Rollback()
	messRepo := uow.MessageRepository()
	msg := &model.Message{ChatID: req.ChatID, UserID: req.SenderID, Message: req.Message}
	err = messRepo.AddMessage(msg)
	if err != nil {
		logger.Error("failed to add message", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	logger.Info("message added", "id", msg.ID, "chat_id", msg.ChatID, "user_id", msg.UserID, "message", msg.Message)
	return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: "Success"}
}
