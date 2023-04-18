-- +goose Up
ALTER TABLE vacancy_info ADD COLUMN location VARCHAR(255);
ALTER TABLE vacancy_info ADD COLUMN company VARCHAR(255);


-- +goose Down
ALTER TABLE vacancy_info DROP COLUMN location;
ALTER TABLE vacancy_info DROP COLUMN company;

