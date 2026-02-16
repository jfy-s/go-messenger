package handlers

import (
	"encoding/json"
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type GetAllUserChatsRequest struct {
	UserId uint64 `validate:"required,min=1"`
}

func HandleGetAllUserChats(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	req := &GetAllUserChatsRequest{UserId: msgPacketRequest.From}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUserChats, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}

	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of websocket", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUserChats, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	defer uow.Rollback()
	chatsRepo := uow.ChatRepository()
	chats, err := chatsRepo.GetAllUserChats(req.UserId)
	if err != nil {
		logger.Error("failed to query all user chats", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUserChats, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit user chats", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUserChats, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	response, err := json.Marshal(chats)
	if err != nil {
		logger.Error("failed to marshal response to GetAllUserChats", "error", err)
		return &model.MessagePacketRequest{MsgType: model.GetAllUserChats, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}

	logger.Info("messages received", "count", len(chats), "user_id", req.UserId)
	return &model.MessagePacketRequest{MsgType: model.GetAllUserChats, From: 0, To: msgPacketRequest.From, Data: response}
}
