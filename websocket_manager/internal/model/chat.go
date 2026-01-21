package model

type Chat struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	CreatorID uint64 `json:"creator_id"`
	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at"`
}

type ChatUsers struct {
	ChatID uint64
	UserID uint64
}
