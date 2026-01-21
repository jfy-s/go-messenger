package handlers

import (
	"log/slog"
	"strconv"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type AddUserToChatRequest struct {
	CreatorID uint64 `validate:"required, min=1"`
	ChatID    uint64 `validate:"required"`
	UserID    uint64 `validate:"required"`
}

func HandleAddUserToChat(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	userId, err := strconv.ParseUint(msgPacketRequest.Data, 10, 64)
	if err != nil {
		logger.Error("failed to parse user id", "error", err)
		return &model.MessagePacketRequest{MsgType: model.AddUserToChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	req := AddUserToChatRequest{CreatorID: msgPacketRequest.From, ChatID: msgPacketRequest.To, UserID: userId}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.AddUserToChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.AddUserToChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	defer uow.Rollback()
	chatRepo := uow.ChatRepository()
	ownerId, err := chatRepo.GetOwnerID(req.ChatID)
	if err != nil {
		logger.Error("failed to get owner id", "error", err)
		return &model.MessagePacketRequest{MsgType: model.AddUserToChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	if ownerId != req.CreatorID {
		logger.Error("user is not owner", "error", err)
		return &model.MessagePacketRequest{MsgType: model.AddUserToChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	chatUsers := &model.ChatUsers{ChatID: req.ChatID, UserID: req.UserID}
	err = chatRepo.AddUserToChat(chatUsers)
	if err != nil {
		logger.Error("failed to add user to chat", "error", err)
		return &model.MessagePacketRequest{MsgType: model.AddUserToChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.AddUserToChat, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	logger.Info("user added to chat", "chat_id", req.ChatID, "user_id", req.UserID)
	return &model.MessagePacketRequest{MsgType: model.AddUserToChat, From: 0, To: msgPacketRequest.From, Data: "Success"}
}
