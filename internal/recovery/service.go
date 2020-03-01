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

	jsonBody, err := json.Marshal(emailDto)
	if err != nil {
		return false, config.ErrUnexpectedError
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", r.cfg.EmailSenderAddr, r.cfg.EmailSenderEndpoints.SendEmail), "application/json", bytes.NewBuffer(jsonBody))

	if err != nil || resp.StatusCode != http.StatusOK {
		return false, config.ErrEmailSendServiceNotWorking
	}
	return true, nil
}
