-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chat_users (
    chat_id BIGINT NOT NULL REFERENCES chat(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (chat_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat_users;
DROP TABLE IF EXISTS chat;
-- +goose StatementEnd
