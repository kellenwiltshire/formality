-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION notify_form_submissions_insert()
RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify(
        'form_submissions_inserts',
        NEW.id::text
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

DROP TRIGGER IF EXISTS form_submissions_insert_trigger on form_submissions;

CREATE TRIGGER form_submissions_insert_trigger
AFTER INSERT ON form_submissions
FOR EACH ROW
EXECUTE FUNCTION notify_form_submissions_insert();

-- +goose Down
DROP TRIGGER IF EXISTS form_submissions_insert_trigger on form_submissions;

DROP FUNCTION IF EXISTS notify_form_submissions_insert();
