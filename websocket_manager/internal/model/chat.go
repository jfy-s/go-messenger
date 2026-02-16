package model

import "time"

type Chat struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatorID uint64    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatUsers struct {
	ChatID uint64
	UserID uint64
}
