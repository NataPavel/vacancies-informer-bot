-- +goose Up
CREATE TABLE error_log
(
    id          SERIAL PRIMARY KEY,
    error_message    TEXT,
    created_at  TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE error_log;

