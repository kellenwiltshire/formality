-- +goose Up
-- +goose StatementBegin
CREATE TYPE status_type AS ENUM ('received', 'dispatched', 'error')
-- +goose StatementEnd

ALTER TABLE form_submissions
    ADD COLUMN status status_type DEFAULT 'received' NOT NULL;

-- +goose Down
-- +goose StatementBegin
ALTER TABLE form_submissions DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS status_type;
-- +goose StatementEnd
