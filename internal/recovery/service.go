package recovery

import (
	"fmt"
	"net/http"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"
	"github.com/CienciaArgentina/go-email-sender/commons"
	"github.com/CienciaArgentina/go-email-sender/defines"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/CienciaArgentina/go-enigma/internal/encryption"
)

const (
	ErrEmailByUserIdFetchCode = "cant_fetch_email"

	// Email verification.
	ErrEmailAlreadyVerified     = "El mail ya se encuentra confirmado"
	ErrEmailAlreadyVerifiedCode = "verified_email"

	// Sending email failed.
	ErrEmailSendingFailed     = "El envío de email falló"
	ErrEmailSendingFailedCode = "failed_email_send"

	// Empty field.
	ErrEmailValidationFailed     = "La validación del email falló por algún campo vacío"
	ErrEmailValidationFailedCode = "empty_field_validating"

	ErrPasswordConfirmationDoesntMatch     = "Los passwords ingresados no son idénticos"
	ErrPasswordConfirmationDoesntMatchCode = "password_mismatch"

	ErrPasswordTokenIsNotValid     = "El token para resetear la contraseña no es válido"
	ErrPasswordTokenIsNotValidCode = "invalid_token"

	errFailedDecryptionCode = "failed_decryption"
)

type recoveryService struct {
	repository RecoveryRepository
	cfg        *config.EnigmaConfig
}

func NewService(cfg *config.EnigmaConfig, r RecoveryRepository) RecoveryService {
	return &recoveryService{
		repository: r,
		cfg:        cfg,
	}
}

func (r *recoveryService) SendConfirmationEmail(userId int64) (bool, apierror.ApiError) {
	verificationToken, userEmail, err := r.repository.GetEmailByUserId(userId)
	if err != nil {
		// TODO: log this
		return false, err
	}

	// If the email or register doesn't exist we should tell the register that an email has been sent IF the email exist. Just to preserve users privacy
	if verificationToken == "" || userEmail == nil || userEmail == (&domain.UserEmail{}) {
		return true, nil
	}

	if userEmail.VerfiedEmail {
		return false, apierror.New(http.StatusBadRequest, ErrEmailAlreadyVerified, apierror.NewErrorCause(ErrEmailAlreadyVerified, ErrEmailAlreadyVerifiedCode))
	}

	url := fmt.Sprintf("%s%s%s?email=%s&token=%s", r.cfg.Microservices.BaseUrl, r.cfg.Microservices.UsersEndpoints.BaseResource, r.cfg.UsersEndpoints.ConfirmEmail,
		userEmail.Email, verificationToken)

	emailDto := commons.NewDTO([]string{userEmail.Email}, url, defines.ConfirmEmail)

	sent, e, _ := rest.EmailSenderApiCall(&r.cfg.Microservices, emailDto)
	if e != nil || !sent {
		return sent, apierror.New(http.StatusBadRequest, ErrEmailSendingFailed, apierror.NewErrorCause(e.Error(), ErrEmailSendingFailedCode))
	}

	return true, nil
}

func (r *recoveryService) ConfirmEmail(email string, token string) (bool, apierror.ApiError) {
	if email == "" || token == "" {
		return false, apierror.New(http.StatusBadRequest, ErrEmailValidationFailed, apierror.NewErrorCause(ErrEmailValidationFailed, ErrEmailValidationFailedCode))
	}

	err := r.repository.ConfirmUserEmail(email, token)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *recoveryService) ResendEmailConfirmationEmail(email string) (bool, apierror.ApiError) {
	if email == "" {
		return false, apierror.New(http.StatusBadRequest, domain.ErrEmptyEmail, apierror.NewErrorCause(domain.ErrEmptyEmail, domain.ErrEmptyEmailCode))
	}

	userId, err := r.repository.GetuserIdByEmail(email)
	if err != nil {
		return false, err
	}

	sent, err := r.SendConfirmationEmail(userId)
	if err != nil || !sent {
		return false, err
	}

	return sent, nil
}

func (r *recoveryService) SendUsername(email string) (bool, apierror.ApiError) {
	if email == "" {
		return false, apierror.New(http.StatusBadRequest, domain.ErrEmptyEmail, apierror.NewErrorCause(domain.ErrEmptyEmail, domain.ErrEmptyEmailCode))
	}

	username, err := r.repository.GetUsernameByEmail(email)
	if err != nil {
		return false, err
	}

	emailDto := commons.NewDTO([]string{email}, username, defines.ForgotUsername)

	sent, e, _ := rest.EmailSenderApiCall(&r.cfg.Microservices, emailDto)
	if e != nil || !sent {
		return sent, apierror.New(http.StatusBadRequest, ErrEmailSendingFailed, apierror.NewErrorCause(e.Error(), ErrEmailSendingFailedCode))
	}

	return sent, nil
}

func (r *recoveryService) SendPasswordReset(email string) (bool, apierror.ApiError) {
	if email == "" {
		return false, apierror.New(http.StatusBadRequest, domain.ErrEmptyEmail, apierror.NewErrorCause(domain.ErrEmptyEmail, domain.ErrEmptyEmailCode))
	}

	securityToken, err := r.repository.GetSecurityToken(email)
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s%s%s?email=%s&token=%s", r.cfg.Microservices.BaseUrl, r.cfg.Microservices.UsersEndpoints.BaseResource, r.cfg.UsersEndpoints.SendPasswordReset,
		email, securityToken)

	emailDto := commons.NewDTO([]string{email}, url, defines.SendPasswordReset)

	sent, e, _ := rest.EmailSenderApiCall(&r.cfg.Microservices, emailDto)
	if e != nil || !sent {
		return sent, apierror.New(http.StatusBadRequest, ErrEmailSendingFailed, apierror.NewErrorCause(e.Error(), ErrEmailSendingFailedCode))
	}

	return sent, nil
}

func (r *recoveryService) ResetPassword(email, password, confirmPassword, token string) (bool, apierror.ApiError) {
	if email == "" || password == "" || confirmPassword == "" || token == "" {
		return false, apierror.New(http.StatusBadRequest, domain.ErrEmptyField, apierror.NewErrorCause(domain.ErrEmptyField, domain.ErrEmptyFieldCode))
	}

	if password != confirmPassword {
		return false, apierror.New(http.StatusBadRequest, ErrPasswordConfirmationDoesntMatch, apierror.NewErrorCause(ErrPasswordConfirmationDoesntMatch,
			ErrPasswordConfirmationDoesntMatchCode))
	}

	securityToken, err := r.repository.GetSecurityToken(email)
	if err != nil {
		return false, err
	}

	if token != securityToken {
		return false, apierror.New(http.StatusBadRequest, ErrPasswordTokenIsNotValid, apierror.NewErrorCause(ErrPasswordTokenIsNotValid, ErrPasswordTokenIsNotValidCode))
	}

	newHashedPassword, e := encryption.GenerateEncodedHash(password, r.cfg)
	if e != nil {
		return false, apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(e.Error(), errFailedDecryptionCode))
	}

	newSecurityToken, errVal := encryption.GenerateSecurityToken(password, r.cfg)
	if errVal != nil {
		return false, apierror.NewInternalServerApiError(errVal.Error(), errVal, "security_token_err")
	}

	userId, err := r.repository.GetuserIdByEmail(email)
	if err != nil {
		return false, err
	}

	updated, err := r.repository.UpdatePasswordHash(userId, newHashedPassword)
	if err != nil {
		return false, err
	}

	if updated {
		emailDto := commons.DTO{
			To:       []string{email},
			Data:     nil,
			Template: "passwordresetnotification",
		}

		_, e, _ := rest.EmailSenderApiCall(&r.cfg.Microservices, &emailDto)
		if e != nil {
			// log this
		}
	}

	_, err = r.repository.UpdateSecurityToken(userId, newSecurityToken)
	if err != nil {
		// TODO: LOG THIS
	}

	return updated, nil
}

func (r *recoveryService) GetUserByUserId(userId int64) (*domain.User, apierror.ApiError) {
	usr, err := r.repository.GetUserByUserId(userId)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
