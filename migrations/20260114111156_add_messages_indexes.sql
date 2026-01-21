-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages USING HASH (chat_id);
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages USING HASH (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_chat_id;
DROP INDEX IF EXISTS idx_messages_user_id;
-- +goose StatementEnd
