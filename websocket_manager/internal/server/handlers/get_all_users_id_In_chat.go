package handlers

import (
	"encoding/json"
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type GetAllUsersIDInChatRequest struct {
	UserID uint64 `validate:"required,min=1"`
	ChatID uint64 `validate:"required"`
}

func HandleGetLlUsersIDInChat(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	req := GetAllUsersIDInChatRequest{UserID: msgPacketRequest.From, ChatID: msgPacketRequest.To}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUsersIDInChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUsersIDInChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	defer uow.Rollback()
	chatRepo := uow.ChatRepository()
	chatUsersIDs := make([]uint64, 0)
	chatUsersIDs, err = chatRepo.GetAllUsersIDInChat(req.ChatID)
	if err != nil {
		logger.Error("failed to get users from chat", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUsersIDInChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUsersIDInChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	logger.Info("users received", "count", len(chatUsersIDs), "chat_id", req.ChatID, "user_id", req.UserID)
	return &model.MessagePacketRequest{MsgType: model.GetAllUsersIDInChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Success")}
}
