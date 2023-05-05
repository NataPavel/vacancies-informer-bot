-- +goose Up
CREATE TABLE user_chat_id
(
    id          SERIAL PRIMARY KEY,
    chat_id     INT,
    created_at  TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE user_chat_id;