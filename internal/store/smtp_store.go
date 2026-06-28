package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

type encryption struct {
	plaintext *string
	cipher    string
}

func (e *encryption) EncryptAES(plaintext string) error {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return fmt.Errorf("Error getting secret key from environment")
	}

	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return fmt.Errorf("invalid key format: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	e.plaintext = &plaintext
	e.cipher = base64.StdEncoding.EncodeToString(ciphertext)

	return nil
}

func (e *encryption) DecryptAES(encrypted string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("Error getting secret key from environment")
	}

	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("invalid key format: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("invalid ciphertext format: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	encryptedData := ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

type Smtp struct {
	Id                int        `json:"id"`
	UserID            int        `json:"user_id"`
	Host              string     `json:"host"`
	Port              int        `json:"port"`
	Username          string     `json:"username"`
	PasswordEncrypted encryption `json:"-"`
	EncryptionType    string     `json:"encryption_type"`
	UpdatedAt         string     `json:"updated_at"`
	RecipientEmail    string     `json:"recipient_email"`
	SenderEmail       string     `json:"sender_email"`
}

type PostgresSmtpStore struct {
	db *sql.DB
}

func NewPostgresSmtpStore(db *sql.DB) *PostgresSmtpStore {
	return &PostgresSmtpStore{
		db: db,
	}
}

type SmtpStore interface {
	CreateSmtpSettings(smtp *Smtp) error
	GetSmtpSettings(userId int64) (*Smtp, error)
	UpdateSmtpSettings(smtp *Smtp) error
	DeleteSmtpSettings(userId int64) error
	GetSmtpEmailSettings(userId int64) (*Smtp, string, error)
}

func (s *PostgresSmtpStore) CreateSmtpSettings(smtp *Smtp) error {
	query := `
		INSERT INTO smtp_Settings (user_id, host, port, username, password_encrypted, encryption_type, recipient_email, sender_email) values ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`

	err := s.db.QueryRow(query, smtp.UserID, smtp.Host, smtp.Port, smtp.Username, smtp.PasswordEncrypted.cipher, smtp.EncryptionType, smtp.RecipientEmail, smtp.SenderEmail).Scan(&smtp.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresSmtpStore) GetSmtpSettings(user_id int64) (*Smtp, error) {
	smtp := &Smtp{}

	query := `
		SELECT id, user_id, host, port, username, encryption_type, updated_at, recipient_email, sender_email FROM smtp_settings WHERE user_id = $1
	`

	err := s.db.QueryRow(query, user_id).Scan(&smtp.Id, &smtp.UserID, &smtp.Host, &smtp.Port, &smtp.Username, &smtp.EncryptionType, &smtp.UpdatedAt, &smtp.RecipientEmail, &smtp.SenderEmail)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return smtp, nil
}

func (s *PostgresSmtpStore) UpdateSmtpSettings(smtp *Smtp) error {
	query := `
		UPDATE smtp_Settings SET host = $1, port = $2, username = $3, password_encrypted = $4, encryption_type = $5, updated_at = CURRENT_TIMESTAMP, recipient_email = $6, sender_email = $7 WHERE id = $8
	`
	result, err := s.db.Exec(query, smtp.Host, smtp.Port, smtp.Username, smtp.PasswordEncrypted.cipher, smtp.EncryptionType, smtp.RecipientEmail, smtp.SenderEmail, smtp.Id)
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

func (s *PostgresSmtpStore) DeleteSmtpSettings(user_id int64) error {
	query := `
		DELETE FROM smtp_settings WHERE user_id = $1
	`

	result, err := s.db.Exec(query, user_id)
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

func (s *PostgresSmtpStore) GetSmtpEmailSettings(user_id int64) (*Smtp, string, error) {
	smtp := &Smtp{}

	query := `
		SELECT id, user_id, host, port, username, password_encrypted, encryption_type, updated_at, recipient_email, sender_email FROM smtp_settings WHERE user_id = $1
	`

	var password string

	err := s.db.QueryRow(query, user_id).Scan(&smtp.Id, &smtp.UserID, &smtp.Host, &smtp.Port, &smtp.Username, &password, &smtp.EncryptionType, &smtp.UpdatedAt, &smtp.RecipientEmail, &smtp.SenderEmail)
	if err == sql.ErrNoRows {
		return nil, "", nil
	}

	if err != nil {
		return nil, "", err
	}

	decrypted_pass, err := smtp.PasswordEncrypted.DecryptAES(password)

	smtp.PasswordEncrypted.plaintext = &decrypted_pass

	return smtp, decrypted_pass, nil
}
