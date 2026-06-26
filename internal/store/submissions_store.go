package store

import "database/sql"

type Submission struct {
	Id          int64  `json:"id"`
	FormId      string `json:"form_id"`
	Payload     string `json:"payload"`
	SubmittedAt string `json:"submitted_at"`
	Status      string `json:"status"`
}

type PostgresSubmissionsStore struct {
	db *sql.DB
}

func NewPostgresSubmissionsStore(db *sql.DB) *PostgresSubmissionsStore {
	return &PostgresSubmissionsStore{
		db: db,
	}
}

type SubmissionsStore interface {
	CreateSubmission(submission *Submission) error
	GetFormSubmissions(form_id string) ([]Submission, error)
	GetFormSubmissionById(submission_id int64) (*Submission, error)
	DeleteSubmission(submission_id int64) error
	DoesUserOwnForm(user_id int64, form_id string) error
	UpdateSubmissionStatus(submission_id int64, status string) error
}

func (s *PostgresSubmissionsStore) CreateSubmission(submission *Submission) error {
	query := `
		INSERT INTO form_submissions (form_id, payload) values ($1, $2) RETURNING id
	`

	err := s.db.QueryRow(query, submission.FormId, submission.Payload).Scan(&submission.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresSubmissionsStore) GetFormSubmissions(form_id string) ([]Submission, error) {
	query := `
		SELECT id, form_id, payload, submitted_at, status FROM form_submissions WHERE form_id = $1
	`

	var submissions []Submission
	rows, err := s.db.Query(query, form_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var submission Submission
		if err := rows.Scan(&submission.Id, &submission.FormId, &submission.Payload, &submission.SubmittedAt, &submission.Status); err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func (s *PostgresSubmissionsStore) GetFormSubmissionById(submission_id int64) (*Submission, error) {
	query := `
		SELECT id, form_id, payload, submitted_at, status FROM form_submissions WHERE id = $1
	`

	submission := &Submission{}
	err := s.db.QueryRow(query, submission_id).Scan(&submission.Id, &submission.FormId, &submission.Payload, &submission.SubmittedAt, &submission.Status)
	if err != nil {
		return nil, err
	}
	return submission, nil

}

func (s *PostgresSubmissionsStore) DeleteSubmission(submission_id int64) error {
	query := `
		DELETE FROM form_submissions WHERE id = $1
	`

	result, err := s.db.Exec(query, submission_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *PostgresSubmissionsStore) DoesUserOwnForm(user_id int64, form_id string) error {
	query := `
		"SELECT COUNT(*) FROM forms WHERE id = $1 AND user_id = $2"
	`

	var count *int64
	err := s.db.QueryRow(query, form_id, user_id).Scan(&count)
	if err != nil || err == sql.ErrNoRows {
		return err
	}

	return nil
}

func (s *PostgresSubmissionsStore) UpdateSubmissionStatus(submission_id int64, status string) error {
	query := `
		UPDATE form_submissions SET status = $1 WHERE id = $2
	`

	result, err := s.db.Exec(query, status, submission_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
