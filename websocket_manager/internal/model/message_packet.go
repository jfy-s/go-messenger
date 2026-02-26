package model

import (
	"encoding/json"
)

type MsgType int

const (
	GetMessage = iota
	SendMessage
	UpdateMessage
	DeleteMessage
	GetAllMessagesInChat // should be limited to some reasonable amount
	CreateChat
	UpdateChat
	DeleteChat
	AddUserToChat
	DeleteUserFromChat
	GetAllUsersIDInChat
	GetAllUserChats
)

const (
	Success       = `"Success"`
	InternalError = `"InternalError"`
)

type MessagePacketRequest struct {
	MsgType MsgType         `json:"msgType"`
	From    uint64          `json:"from,omitempty"`
	To      uint64          `json:"to,omitempty"`
	Data    json.RawMessage `json:"data"`
}

func ByteToMessagePacketRequest(b []byte) (*MessagePacketRequest, error) {
	var msgPkt MessagePacketRequest
	err := json.Unmarshal(b, &msgPkt)
	return &msgPkt, err
}

func (msg *MessagePacketRequest) ToBytes() ([]byte, error) {
	return json.Marshal(msg)
}
