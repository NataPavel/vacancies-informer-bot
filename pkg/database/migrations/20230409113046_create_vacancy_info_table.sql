-- +goose Up
CREATE TABLE vacancy_info
(
    id          SERIAL PRIMARY KEY,
    website     VARCHAR(255),
    vacancy_link TEXT,
    vacancy_text TEXT,
    link_hash    TEXT,
    created_at  TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE vacancy_info;

