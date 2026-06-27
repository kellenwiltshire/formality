package store

import (
	"database/sql"
)

type Form struct {
	Id          string `json:"id"`
	UserId      int    `json:"user_id"`
	Name        string `json:"name"`
	TargetEmail string `json:"target_email"`
	CreatedAt   string `json:"created_at"`
}

type PostgresFormStore struct {
	db *sql.DB
}

func NewPostgresFormStore(db *sql.DB) *PostgresFormStore {
	return &PostgresFormStore{
		db: db,
	}
}

type FormStore interface {
	CreateForm(*Form) error
	GetForm(form_id string, user_id int64) (*Form, error)
	UpdateForm(*Form) error
	DeleteForm(form_id string, user_id int64) error
	GetAllFormsForUser(user_id int64) ([]Form, error)
	GetFormInfoForEmail(form_id string) (*Form, error)
}

func (s *PostgresFormStore) CreateForm(form *Form) error {
	query := `
		INSERT INTO forms (user_id, name, target_email) values ($1, $2, $3) RETURNING id
	`

	err := s.db.QueryRow(query, form.UserId, form.Name, form.TargetEmail).Scan(&form.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresFormStore) GetForm(form_id string, user_id int64) (*Form, error) {
	form := &Form{}

	query := `
		SELECT id, user_id, name, target_email, created_at FROM forms WHERE id = $1 AND user_id = $2
	`

	err := s.db.QueryRow(query, form_id, user_id).Scan(&form.Id, &form.UserId, &form.Name, &form.TargetEmail, &form.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return form, nil

}

func (s *PostgresFormStore) UpdateForm(form *Form) error {
	query := `
		UPDATE forms SET name = $1, target_email = $2 WHERE id = $3 AND user_id = $4
	`

	result, err := s.db.Exec(query, form.Name, form.TargetEmail, form.Id, form.UserId)
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

func (s *PostgresFormStore) DeleteForm(form_id string, user_id int64) error {
	query := `
		DELETE FROM forms WHERE id = $1 AND user_id = $2
	`

	result, err := s.db.Exec(query, form_id, user_id)
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

func (s *PostgresFormStore) GetAllFormsForUser(user_id int64) ([]Form, error) {
	query := `
		SELECT id, user_id, name, target_email, created_at FROM forms WHERE user_id = $1
	`

	rows, err := s.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}

	var forms []Form
	for rows.Next() {
		var form Form
		err := rows.Scan(&form.Id, &form.UserId, &form.Name, &form.TargetEmail, &form.CreatedAt)
		if err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}
	return forms, nil
}

func (s *PostgresFormStore) GetFormInfoForEmail(form_id string) (*Form, error) {
	query := `
		SELECT user_id, target_email, name FROM forms WHERE id = $1
	`

	form := &Form{}

	err := s.db.QueryRow(query, form_id).Scan(&form.UserId, &form.TargetEmail, &form.Name)
	if err != nil {
		return nil, err
	}

	return form, nil
}
