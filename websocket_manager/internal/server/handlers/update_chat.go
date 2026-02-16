package handlers

import (
	"encoding/json"
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type UpdateChatRequest struct {
	CreatorID uint64 `validate:"required,min=1"`
	ChatID    uint64 `validate:"required"`
	Name      string `validate:"required,min=1,max=64"`
}

func HandleUpdateChat(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	var name string
	_ = json.Unmarshal(msgPacketRequest.Data, &name)
	req := UpdateChatRequest{CreatorID: msgPacketRequest.From, ChatID: msgPacketRequest.To, Name: name}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.UpdateChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.UpdateChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	defer uow.Rollback()
	chatRepo := uow.ChatRepository()
	ownerId, err := chatRepo.GetOwnerID(req.ChatID)
	if err != nil {
		logger.Error("failed to get owner id", "error", err)
		return &model.MessagePacketRequest{MsgType: model.UpdateChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	if ownerId != req.CreatorID {
		logger.Error("user is not owner", "error", err)
		return &model.MessagePacketRequest{MsgType: model.UpdateChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	chat := &model.Chat{CreatorID: req.CreatorID, ID: req.ChatID, Name: req.Name}
	err = chatRepo.UpdateChat(chat)
	if err != nil {
		logger.Error("failed to update chat", "error", err)
		return &model.MessagePacketRequest{MsgType: model.UpdateChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.UpdateChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	logger.Info("chat updated", "chat_id", chat.ID, "name", chat.Name)
	return &model.MessagePacketRequest{MsgType: model.UpdateChat, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Success")}
}
