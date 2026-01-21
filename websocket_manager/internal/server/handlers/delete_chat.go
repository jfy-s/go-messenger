package handlers

import (
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type DeleteChatRequest struct {
	CreatorID uint64 `validate:"required, min=1"`
	ChatID    uint64 `validate:"required"`
}

func HandleDeleteChat(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	req := DeleteChatRequest{CreatorID: msgPacketRequest.From, ChatID: msgPacketRequest.To}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	defer uow.Rollback()
	chatRepo := uow.ChatRepository()
	ownerId, err := chatRepo.GetOwnerID(req.ChatID)
	if err != nil {
		logger.Error("failed to get owner id", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	if ownerId != req.CreatorID {
		logger.Error("user is not owner", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	err = chatRepo.DeleteChat(req.ChatID)
	if err != nil {
		logger.Error("failed to delete chat", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	logger.Info("chat deleted", "chat_id", req.ChatID)
	return &model.MessagePacketRequest{MsgType: model.DeleteChat, From: 0, To: msgPacketRequest.From, Data: "Success"}
}
