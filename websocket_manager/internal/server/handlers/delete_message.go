package handlers

import (
	"log/slog"
	"strconv"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"

	"github.com/go-playground/validator"
)

type DeleteMessageRequest struct {
	DeletterID uint64 `validate:"required, min=1"`
	ChatID     uint64 `validate:"required"`
	MsgID      uint64 `validate:"required"`
}

func HandleDeleteMessage(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	msgID, err := strconv.ParseUint(msgPacketRequest.Data, 10, 64)
	if err != nil {
		logger.Error("failed to parse message id", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	req := DeleteMessageRequest{DeletterID: msgPacketRequest.From, ChatID: msgPacketRequest.To, MsgID: msgID}
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		logger.Error("failed to validate request", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	defer uow.Rollback()
	msgRepo := uow.MessageRepository()
	msgSender, err := msgRepo.GetSenderID(req.MsgID)
	if err != nil {
		logger.Error("failed to get message sender", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	chatRepo := uow.ChatRepository()
	ownerID, err := chatRepo.GetOwnerID(req.ChatID)
	if err != nil {
		logger.Error("failed to get owner id", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	if msgSender != req.DeletterID || ownerID != req.DeletterID {
		logger.Error("user is not message sender or chat owner", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	err = msgRepo.DeleteMessage(req.MsgID)
	if err != nil {
		logger.Error("failed to delete message", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	err = uow.Commit()
	if err != nil {
		logger.Error("failed to commit unit of work", "error", err)
		return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Internal Error"}
	}
	logger.Info("message deleted", "id", req.MsgID, "deleted_by", req.DeletterID)
	return &model.MessagePacketRequest{MsgType: model.DeleteMessage, From: 0, To: msgPacketRequest.From, Data: "Success"}
}
