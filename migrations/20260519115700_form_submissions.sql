-- +goose Up
CREATE TABLE IF NOT EXISTS form_submissions (
    id SERIAL PRIMARY KEY,
    form_id VARCHAR(6) NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    payload JSONB NOT NULL,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_submissions_form_id ON form_submissions(form_id);

-- +goose Down
DROP TABLE IF EXISTS form_submissions;