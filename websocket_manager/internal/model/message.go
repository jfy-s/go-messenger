package model

import "time"

type Message struct {
	ID        uint64    `json:"id"`
	ChatID    uint64    `json:"chat_id"`
	UserID    uint64    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MessageToByte(m *Message) []byte {
	return []byte(m.Message)
}

func ByteToMessage(b []byte) *Message {
	return &Message{Message: string(b)}
}
