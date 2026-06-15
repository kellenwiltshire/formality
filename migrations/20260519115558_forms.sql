-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_short_id() 
RETURNS TEXT AS $$
DECLARE
    chars TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    result TEXT := '';
    i INT;
BEGIN
    FOR i IN 1..6 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql VOLATILE;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS forms (
    id VARCHAR(6) PRIMARY KEY DEFAULT generate_short_id(), 
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    target_email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_forms_user_id ON forms(user_id);

-- +goose Down
DROP TABLE IF EXISTS forms;
DROP FUNCTION IF EXISTS generate_short_id();