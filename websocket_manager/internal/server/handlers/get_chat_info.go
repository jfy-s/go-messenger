package handlers

import (
	"encoding/json"
	"log/slog"
	"websocket_manager/internal/model"
	"websocket_manager/internal/storage"
)

type GetChatInfoResponse struct {
	ChatInfo *model.Chat  `json:"chat_info"`
	Users    []model.User `json:"users"`
}

func HandleGetChatInfo(storage storage.Storage, msgPacketRequest *model.MessagePacketRequest, logger *slog.Logger) *model.MessagePacketRequest {
	uow, err := storage.CreateUnitOfWork()
	if err != nil {
		logger.Error("failed to create unit of websocket manager", "error", err)
		return &model.MessagePacketRequest{}
	}
	defer uow.Rollback()

	chatInfo, users, err := uow.ChatRepository().GetChatInfo(msgPacketRequest.To)
	if err != nil {
		logger.Error("failed to get chat info", "error", err)
		return &model.MessagePacketRequest{}
	}

	rawResponse := GetChatInfoResponse{ChatInfo: chatInfo, Users: users}
	response, err := json.Marshal(rawResponse)
	if err != nil {
		logger.Error("failed to marshal chat info", "error", err)
		return &model.MessagePacketRequest{}
	}
	logger.Info("chat information")
	return &model.MessagePacketRequest{MsgType: model.GetChatInfo, From: msgPacketRequest.From, To: msgPacketRequest.To, Data: response}
}
