-- +goose Up
-- +goose StatementBegin
CREATE TYPE status_type AS ENUM ('received', 'dispatched', 'error')
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS form_submissions (
    id SERIAL PRIMARY KEY,
    form_id VARCHAR(6) NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    payload JSONB NOT NULL,
    status status_type DEFAULT 'received' NOT NULL,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_submissions_form_id ON form_submissions(form_id);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS form_submissions;
DROP TYPE IF EXISTS status_type;
-- +goose StatementEnd