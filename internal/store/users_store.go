package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/kellenwiltshire/formality/internal/tokens"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintText *string
	hash       []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintText = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err //internal server error
		}
	}

	return true, nil
}

type User struct {
	Id           int      `json:"id"`
	Email        string   `json:"email"`
	PasswordHash password `json:"-"`
	Role         string   `json:"role"`
	CreatedAt    string   `json:"created_at"`
}

var AnonymousUser = &User{} // EVERYONE WHO IS NOT LOGGED IN

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUser(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(*User) error
	DeleteUser(id int64) error
	GetAllUsers() ([]User, error)
	GetNumberAdminUsers() (*int64, error)
	GetUserToken(plaintextPassword string) (*User, error)
	GetAdminToken(plaintextPassword string) (*User, error)
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (email, password_hash, role) values ($1, $2, $3) RETURNING id
	`

	err := s.db.QueryRow(query, user.Email, user.PasswordHash.hash, user.Role).Scan(&user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresUserStore) GetUser(id int64) (*User, error) {
	user := &User{}

	query := `
		SELECT id, email, role, created_at FROM users WHERE id = $1
	`

	err := s.db.QueryRow(query, id).Scan(&user.Id, &user.Email, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) GetUserByEmail(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
	SELECT id, email, password_hash, created_at, role
	FROM users
	WHERE email = $1
	`

	err := s.db.QueryRow(query, username).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash.hash,
		&user.CreatedAt,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users SET email = $1, role = $2 WHERE id = $3
	`
	result, err := s.db.Exec(query, user.Email, user.Role, user.Id)
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

func (s *PostgresUserStore) DeleteUser(id int64) error {
	query := `
		DELETE FROM users WHERE id = $1
	`

	result, err := s.db.Exec(query, id)
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

func (s *PostgresUserStore) GetAllUsers() ([]User, error) {
	query := `
		SELECT id, email, role, created_at FROM users
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *PostgresUserStore) GetNumberAdminUsers() (*int64, error) {
	query := `
		SELECT COUNT(*) FROM users WHERE role = $1
	`
	var count *int64
	err := s.db.QueryRow(query, tokens.ScopeAdmin).Scan(&count)
	if err != nil {
		return nil, err
	}

	return count, nil
}

func (s *PostgresUserStore) GetUserToken(token_hash string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(token_hash))

	query := `
		SELECT u.id, u.email, u.created_at, u.role
		FROM users u
		INNER JOIN tokens t ON t.user_id = u.id
		WHERE t.hash = $1 AND t.expiry > $2
  	`

	user := &User{}

	err := s.db.QueryRow(query, tokenHash[:], time.Now()).Scan(
		&user.Id,
		&user.Email,
		&user.CreatedAt,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) GetAdminToken(token_hash string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(token_hash))

	query := `
		SELECT u.id, u.email, u.created_at, u.role
		FROM users u
		INNER JOIN tokens t ON t.user_id = u.id
		WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3
  	`

	user := &User{}

	err := s.db.QueryRow(query, tokenHash[:], tokens.ScopeAdmin, time.Now()).Scan(
		&user.Id,
		&user.Email,
		&user.CreatedAt,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
