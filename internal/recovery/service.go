package recovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CienciaArgentina/go-email-sender/commons"
	"github.com/CienciaArgentina/go-email-sender/defines"
	"github.com/CienciaArgentina/go-enigma/config"
	"net/http"
)

type Service interface {
	SendConfirmationEmail(userId int64) (bool, error)
	ConfirmEmail(email string, token string) (bool, error)
	ResendEmailConfirmationEmail(email string) (bool, error)
	SendUsername(email string) (bool, error)
	SendRecoveryEmail(dto *commons.DTO) (bool, error)
}

type recoveryService struct {
	repo Repository
	cfg  *config.Configuration
}

func NewService(r Repository, c *config.Configuration) Service {
	return &recoveryService{repo: r, cfg: c}
}

func (r *recoveryService) SendConfirmationEmail(userId int64) (bool, error) {
	if userId == 0 {
		return false, config.ErrEmptyUserId
	}

	verificationToken, userEmail, err := r.repo.GetEmailByUserId(userId)
	if err != nil {
		// TODO: log this
		return false, config.ErrUnexpectedError
	}

	// If the email or user doesn't exist we should tell the user that an email has been sent IF the email exist. Just to preserve users privacy
	if verificationToken == "" || userEmail == nil || userEmail == (&UserEmail{}) {
		return true, nil
	}

	if userEmail.VerfiedEmail {
		return false, config.ErrEmailAlreadyVerified
	}

	url := fmt.Sprintf("%s%s%s?email=%s&token=%s", r.cfg.Microservices.BaseUrl, r.cfg.Microservices.UsersEndpoints.BaseResource, r.cfg.UsersEndpoints.ConfirmEmail,
		userEmail.Email, verificationToken)

	emailDto := commons.NewDTO([]string{userEmail.Email}, url, defines.ConfirmEmail)

	sent, err := r.SendRecoveryEmail(emailDto)
	if err != nil || !sent {
		return sent, err
	}

	return true, nil
}

func (r *recoveryService) ConfirmEmail(email string, token string) (bool, error) {
	if email == "" || token == "" {
		return false, config.ErrEmailValidationFailed
	}

	err := r.repo.ConfirmUserEmail(email, token)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *recoveryService) ResendEmailConfirmationEmail(email string) (bool, error) {
	if email == "" {
		return false, config.ErrEmptyEmail
	}

	userId, err := r.repo.GetuserIdByEmail(email)
	if err != nil {
		return false, err
	}

	sent, err := r.SendConfirmationEmail(userId)
	if err != nil || !sent {
		return false, err
	}

	return sent, nil
}

func (r *recoveryService) SendUsername(email string) (bool, error) {
	if email == "" {
		return false, config.ErrEmptyEmail
	}

	username, err := r.repo.GetUsernameByEmail(email)
	if err != nil {
		return false, err
	}

	emailDto := commons.NewDTO([]string{email}, username, defines.ForgotUsername)

	sent, err := r.SendRecoveryEmail(emailDto)
	if err != nil || !sent {
		return sent, err
	}

	return sent, nil
}

func (r *recoveryService) SendRecoveryEmail(dto *commons.DTO) (bool, error) {
	jsonBody, err := json.Marshal(dto)
	if err != nil {
		return false, config.ErrUnexpectedError
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", r.cfg.EmailSenderAddr, r.cfg.EmailSenderEndpoints.SendEmail), "application/json", bytes.NewBuffer(jsonBody))

	if err != nil || resp.StatusCode != http.StatusOK {
		return false, config.ErrEmailSendServiceNotWorking
	}
	return true, nil
}
