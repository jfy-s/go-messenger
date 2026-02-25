package handlers

import (
	"encoding/json"
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type UpdateMessageRequest struct {
	SenderID uint64 `validate:"required,min=1"`
	MsgID    uint64 `validate:"required"`
	Message  string `validate:"required"`
}

func HandleUpdateMessage(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	var message string
	_ = json.Unmarshal(msgPacketRequest.Data, &message)
	req := UpdateMessageRequest{SenderID: msgPacketRequest.From, MsgID: msgPacketRequest.To, Message: message}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	defer uow.Rollback()
	msgRepo := uow.MessageRepository()
	msgSender, err := msgRepo.GetSenderID(req.SenderID)
	if err != nil {
		logger.Error("failed to get message sender", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	if msgSender != req.SenderID {
		logger.Error("user is not message sender", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	msg := &model.Message{ID: req.MsgID, Message: req.Message}
	err = msgRepo.UpdateMessage(msg)
	if err != nil {
		logger.Error("failed to add message", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Internal Error")}
	}
	logger.Info("message updated", "id", msg.ID, "chat_id", msg.ChatID, "user_id", msg.UserID, "message", msg.Message)
	return &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msgPacketRequest.From, Data: json.RawMessage("Success")}
}
