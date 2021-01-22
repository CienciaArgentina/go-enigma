package recovery

import (
	"errors"
	"fmt"
	"github.com/CienciaArgentina/go-backend-commons/pkg/middleware"
	"net/http"
	"reflect"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/clog"
	"github.com/CienciaArgentina/go-backend-commons/pkg/performance"
	"github.com/go-resty/resty/v2"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
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


	// Empty field.
	ErrEmailValidationFailed     = "La validación del email falló por algún campo vacío"

	ErrPasswordConfirmationDoesntMatch     = "Los passwords ingresados no son idénticos"

	ErrPasswordTokenIsNotValid     = "El token para resetear la contraseña no es válido"

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

func (r *recoveryService) SendConfirmationEmail(userId int64, ctx *middleware.ContextInformation) (bool, apierror.ApiError) {
	var verificationToken string
	var userEmail *domain.UserEmail
	var err apierror.ApiError
	performance.TrackTime(time.Now(), "GetEmailByUserId", ctx, func() {
		verificationToken, userEmail, err = r.repository.GetEmailByUserId(userId)
	})
	if err != nil {
		clog.Error("Can't send confirmation email", "send-confirmation-email", err, map[string]string{"auth_id": fmt.Sprintf("%d", userId), clog.Subtype: "get-email-by-user-id"})
		return false, err
	}

	// If the email or register doesn't exist we should tell the register that an email has been sent IF the email exist. Just to preserve users privacy
	if verificationToken == "" || userEmail == nil || reflect.DeepEqual(userEmail, &domain.UserEmail{}) {
		return true, nil
	}

	if userEmail.VerfiedEmail {
		return false, apierror.NewBadRequestApiError(ErrEmailAlreadyVerified)
	}

	url := fmt.Sprintf("/confirmemail?email=%s&token=%s", userEmail.Email, verificationToken)

	emailDto := commons.NewDTO([]string{userEmail.Email}, url, defines.ConfirmEmail)

	var response *resty.Response
	var apierr error
	// TODO: Move this to a client
	performance.TrackTime(time.Now(), "EmailSendAPICall", ctx, func() {
		response, apierr = resty.New().SetHostURL(domain.GetEmailSenderBaseURL()).R().SetBody(emailDto).Post("/email")
	})

	if apierr != nil {
		clog.Error("Rest client err", "send-confirmation-email", apierr, map[string]string{"auth_id": fmt.Sprintf("%d", userId), clog.Subtype: "confirm-email-api-call"})
		return false, apierror.NewInternalServerApiError(apierr.Error(), apierr, "cannot_email")
	}

	if response.IsError() {
		clog.Error("Email sender status err", "send-confirmation-email", apierr, map[string]string{"status": response.Status(), clog.Subtype: "confirm-email-api-call", "auth_id": fmt.Sprintf("%d", userId)})
		return false, apierror.NewInternalServerApiError("cant send email", errors.New("cant send email"), "cannot_email")
	}

	return true, nil
}

func (r *recoveryService) ConfirmEmail(email string, token string, ctx *middleware.ContextInformation) (bool, apierror.ApiError) {
	if email == "" || token == "" {
		return false, apierror.NewBadRequestApiError(ErrEmailValidationFailed)
	}

	var err apierror.ApiError
	performance.TrackTime(time.Now(), "ConfirmUserEmail", ctx, func() {
		err = r.repository.ConfirmUserEmail(email, token)
	})

	if err != nil {
		clog.Error("ConfirmUserEmail error", "confirm-email", err, map[string]string{clog.Subtype: "confirm-user-email", "email": email})
		return false, err
	}

	return true, nil
}

func (r *recoveryService) ResendEmailConfirmationEmail(email string, ctx *middleware.ContextInformation) (bool, apierror.ApiError) {
	if email == "" {
		return false, apierror.NewBadRequestApiError(domain.ErrEmptyEmail)
	}

	var userId int64
	var err apierror.ApiError
	performance.TrackTime(time.Now(), "GetuserIdByEmail", ctx, func() {
		userId, err = r.repository.GetuserIdByEmail(email)
	})

	if err != nil {
		return false, err
	}

	var sent bool
	performance.TrackTime(time.Now(), "SendConfirmationEmail", ctx, func() {
		sent, err = r.SendConfirmationEmail(userId, ctx)
	})

	if err != nil || !sent {
		return false, err
	}

	return sent, nil
}

func (r *recoveryService) SendUsername(email string, ctx *middleware.ContextInformation) (bool, apierror.ApiError) {
	if email == "" {
		return false, apierror.NewBadRequestApiError(domain.ErrEmptyEmail)
	}

	var username string
	var err apierror.ApiError
	performance.TrackTime(time.Now(), "GetUsernameByEmail", ctx, func() {
		username, err = r.repository.GetUsernameByEmail(email)
	})

	if err != nil {
		return false, err
	}

	emailDto := commons.NewDTO([]string{email}, username, defines.ForgotUsername)

	var response *resty.Response
	var apierr error
	performance.TrackTime(time.Now(), "SendEmailAPICall", ctx, func() {
		response, apierr = resty.New().SetHostURL(domain.GetEmailSenderBaseURL()).R().SetBody(emailDto).Post("/email")
	})

	if apierr != nil {
		return false, apierror.NewInternalServerApiError(apierr.Error(), apierr, "cannot_email")
	}

	if response.IsError() {
		return false, apierror.NewInternalServerApiError("cant send email", errors.New("cant send email"), "cannot_email")
	}

	return true, nil
}

func (r *recoveryService) SendPasswordReset(email string, ctx *middleware.ContextInformation) (bool, apierror.ApiError) {
	if email == "" {
		return false, apierror.NewBadRequestApiError(domain.ErrEmptyEmail)
	}

	var securityToken string
	var err apierror.ApiError
	performance.TrackTime(time.Now(), "GetSecurityToken", ctx, func() {
		securityToken, err = r.repository.GetSecurityToken(email)
	})

	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("/sendpasswordreset?email=%s&token=%s", email, securityToken)

	emailDto := commons.NewDTO([]string{email}, url, defines.SendPasswordReset)

	var response *resty.Response
	var apierr error
	// TODO: Move this to a client
	performance.TrackTime(time.Now(), "SendEmailAPICall", ctx, func() {
		response, apierr = resty.New().SetHostURL(domain.GetEmailSenderBaseURL()).R().SetBody(emailDto).Post("/email")
	})

	if apierr != nil {
		return false, apierror.NewInternalServerApiError(apierr.Error(), apierr, "cannot_email")
	}

	if response.IsError() {
		return false, apierror.NewInternalServerApiError("cant send email", errors.New("cant send email"), "cannot_email")
	}

	return true, nil
}

func (r *recoveryService) ResetPassword(email, password, confirmPassword, token string, ctx *middleware.ContextInformation) (bool, apierror.ApiError) {
	if email == "" || password == "" || confirmPassword == "" || token == "" {
		return false, apierror.NewBadRequestApiError(domain.ErrEmptyField)
	}

	if password != confirmPassword {
		return false, apierror.NewBadRequestApiError(ErrPasswordConfirmationDoesntMatch)
	}

	var securityToken string
	var err apierror.ApiError
	performance.TrackTime(time.Now(), "GetSecurityToken", ctx, func() {
		securityToken, err = r.repository.GetSecurityToken(email)
	})

	if err != nil {
		return false, err
	}

	if token != securityToken {
		return false, apierror.NewBadRequestApiError(ErrPasswordTokenIsNotValid)
	}

	var newHashedPassword string
	var e error
	performance.TrackTime(time.Now(), "GenerateEncodedHash", ctx, func() {
		newHashedPassword, e = encryption.GenerateEncodedHash(password, r.cfg)
	})

	if e != nil {
		return false, apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(e.Error(), errFailedDecryptionCode))
	}

	var newSecurityToken string
	performance.TrackTime(time.Now(), "GenerateSecurityToken", ctx, func() {
		newSecurityToken, e = encryption.GenerateSecurityToken(password, r.cfg)
	})

	if e != nil {
		return false, apierror.NewInternalServerApiError(e.Error(), e, "security_token_err")
	}

	var userId int64
	performance.TrackTime(time.Now(), "GetuserIdByEmail", ctx, func() {
		userId, err = r.repository.GetuserIdByEmail(email)
	})

	if err != nil {
		return false, err
	}

	var updated bool
	performance.TrackTime(time.Now(), "UpdatePasswordHash", ctx, func() {
		updated, err = r.repository.UpdatePasswordHash(userId, newHashedPassword)
	})

	if err != nil {
		return false, err
	}

	if updated {
		emailDto := commons.DTO{
			To:       []string{email},
			Data:     nil,
			Template: "passwordresetnotification",
		}

		var response *resty.Response
		// TODO: Move this to a client
		performance.TrackTime(time.Now(), "SendEmailAPICall", ctx, func() {
			response, e = resty.New().SetHostURL(domain.GetEmailSenderBaseURL()).R().SetBody(emailDto).Post("/email")
		})

		if e != nil {
			return false, apierror.NewInternalServerApiError(e.Error(), e, "cannot_email")
		}

		if response.IsError() {
			return false, apierror.NewInternalServerApiError("cant send email", errors.New("cant send email"), "cannot_email")
		}
	}

	performance.TrackTime(time.Now(), "UpdateSecurityToken", ctx, func() {
		_, err = r.repository.UpdateSecurityToken(userId, newSecurityToken)
	})

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
