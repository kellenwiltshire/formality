package service

import (
	"formality/internal/store"
	"log"
	netSmtp "net/smtp"
	"strconv"
)

type SendMailService struct {
	formStore        store.FormStore
	submissionsStore store.SubmissionsStore
	smtpStore        store.SmtpStore
	logger           *log.Logger
}

func NewSendMailService(formStore store.FormStore, submissionsStore store.SubmissionsStore, smtpStore store.SmtpStore, logger *log.Logger) *SendMailService {
	return &SendMailService{
		formStore:        formStore,
		submissionsStore: submissionsStore,
		smtpStore:        smtpStore,
		logger:           logger,
	}
}

func (s *SendMailService) SendMail(submission_id string) error {
	submissionId, err := strconv.ParseInt(submission_id, 10, 64)
	if err != nil {
		return err
	}
	submission, err := s.submissionsStore.GetFormSubmissionById(submissionId)
	if err != nil {
		return err
	}
	formId, err := strconv.ParseInt(submission.FormId, 10, 64)
	if err != nil {
		return err
	}
	form, err := s.formStore.GetFormInfoForEmail(formId)
	if err != nil {
		return err
	}

	smtp, pass, err := s.smtpStore.GetSmtpEmailSettings(int64(form.UserId))
	if err != nil {
		return err
	}

	// use exported Plaintext field of PasswordEncrypted
	auth := netSmtp.PlainAuth("", smtp.Username, pass, smtp.Host)

	to := []string{smtp.RecipientEmail}
	msg := []byte("To: " + smtp.RecipientEmail + "\r\n" +
		"Subject: New Form Response For " + form.Name + " From Formality\r\n" +
		"\r\n" +
		submission.Payload)

	err = netSmtp.SendMail(smtp.Host+":"+strconv.Itoa(smtp.Port), auth, smtp.SenderEmail, to, msg)
	if err != nil {
		err = s.submissionsStore.UpdateSubmissionStatus(formId, "error")
		return err
	}

	err = s.submissionsStore.UpdateSubmissionStatus(formId, "dispatched")
	if err != nil {
		return err
	}
	return nil
}

func (s *SendMailService) TestSendMail(userId int64, testPayload string) error {
	smtp, pass, err := s.smtpStore.GetSmtpEmailSettings(int64(userId))
	if err != nil {
		return err
	}

	// use exported Plaintext field of PasswordEncrypted
	auth := netSmtp.PlainAuth("", smtp.Username, pass, smtp.Host)

	to := []string{smtp.RecipientEmail}
	msg := []byte("To: " + smtp.RecipientEmail + "\r\n" +
		"Subject: Testing Email Settings From Formality\r\n" +
		"\r\n" +
		testPayload)

	err = netSmtp.SendMail(smtp.Host+":"+strconv.Itoa(smtp.Port), auth, smtp.SenderEmail, to, msg)
	if err != nil {
		return err
	}

	return nil
}
