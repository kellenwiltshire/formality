-- +goose Up
ALTER TABLE smtp_settings
    ADD COLUMN recipient_email VARCHAR(255) NOT NULL,
    ADD COLUMN sender_email VARCHAR(255) NOT NULL;
-- +goose Down 
ALTER TABLE smtp_settings DROP COLUMN IF EXISTS recipient_email,
    DROP COLUMN IF EXISTS sender_email;