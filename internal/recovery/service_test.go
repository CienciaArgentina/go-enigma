package recovery

import (
	"errors"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	GetEmailByUserId   = "GetEmailByUserId"
	ConfirmUserEmail   = "ConfirmUserEmail"
	GetuserIdByEmail   = "GetuserIdByEmail"
	GetUsernameByEmail = "GetUsernameByEmail"
)

type RecoveryRepositoryMock struct {
	mock.Mock
}

func (r *RecoveryRepositoryMock) GetEmailByUserId(userId int64) (string, *UserEmail, error) {
	args := r.Called(userId)
	return args.Get(0).(string), args.Get(1).(*UserEmail), args.Error(2)
}

func (r *RecoveryRepositoryMock) ConfirmUserEmail(email string, token string) error {
	args := r.Called(email, token)
	return args.Error(0)
}

func (r *RecoveryRepositoryMock) GetuserIdByEmail(email string) (int64, error) {
	args := r.Called(email)
	return args.Get(0).(int64), args.Error(1)
}

func (r *RecoveryRepositoryMock) GetUsernameByEmail(email string) (string, error) {
	args := r.Called(email)
	return args.Get(0).(string), args.Error(1)
}

func GetServiceAndMock() (Service, *RecoveryRepositoryMock) {
	mock := new(RecoveryRepositoryMock)
	svc := NewService(mock, config.New())
	return svc, mock
}

func TestNewServiceShouldReturnANewService(t *testing.T) {
	svc, _ := GetServiceAndMock()
	require.NotNil(t, svc)
}

func TestSendConfirmationEmailShouldReturnEmptyUserIdError(t *testing.T) {
	svc, _ := GetServiceAndMock()
	sent, err := svc.SendConfirmationEmail(0)
	require.False(t, sent)
	require.Equal(t, config.ErrEmptyUserId, err)
}

func TestSendConfirmationEmailShouldReturnErrorWhileGettingEmail(t *testing.T) {
	svc, mock := GetServiceAndMock()
	mock.On(GetEmailByUserId, int64(1)).Return("", &UserEmail{}, errors.New("Indistinct"))
	sent, err := svc.SendConfirmationEmail(1)
	require.Equal(t, config.ErrUnexpectedError, err)
	require.False(t, sent)
}

func TestSendConfirmationEmailShouldReturnNoErrorIfUserEmailIsEmpty(t *testing.T) {
	svc, mock := GetServiceAndMock()
	mock.On(GetEmailByUserId, int64(1)).Return("", &UserEmail{}, nil)
	sent, err := svc.SendConfirmationEmail(1)
	require.NoError(t, err)
	require.True(t, sent)
}

func TestSendConfirmationEmailShouldReturnErrorIfEmailIsAlreadyVerified(t *testing.T) {
	svc, mock := GetServiceAndMock()
	verified := &UserEmail{
		VerfiedEmail: true,
	}
	mock.On(GetEmailByUserId, int64(1)).Return("notempty", verified, nil)
	sent, err := svc.SendConfirmationEmail(1)
	require.Equal(t, config.ErrEmailAlreadyVerified, err)
	require.False(t, sent)
}

func TestConfirmEmailShouldThrowErrorWhenEmailIsEmpty(t *testing.T) {
	svc, _ := GetServiceAndMock()

	confirm, err := svc.ConfirmEmail("", "")
	require.False(t, confirm)
	require.Equal(t, config.ErrEmailValidationFailed, err)
}

func TestConfirmEmailShouldThrowErrorWhenTokenIsNil(t *testing.T) {
	svc, _ := GetServiceAndMock()

	confirm, err := svc.ConfirmEmail("asd", "")
	require.False(t, confirm)
	require.Equal(t, config.ErrEmailValidationFailed, err)
}

func TestConfirmEmailShouldThrowErrorWhenConfirmationFails(t *testing.T) {
	svc, mock := GetServiceAndMock()
	mock.On(ConfirmUserEmail, "asd", "asd").Return(config.ErrEmailValidationFailed)
	confirm, err := svc.ConfirmEmail("asd", "asd")
	require.False(t, confirm)
	require.Equal(t, config.ErrEmailValidationFailed, err)
}

func TestConfirmEmailShouldThrowNoError(t *testing.T) {
	svc, mock := GetServiceAndMock()
	mock.On(ConfirmUserEmail, "asd", "asd").Return(nil)
	confirm, err := svc.ConfirmEmail("asd", "asd")
	require.True(t, confirm)
	require.Nil(t, err)
}

func TestResendEmailConfirmationEmailShouldReturnErrorWhenEmailisNil(t *testing.T) {
	svc, _ := GetServiceAndMock()
	sent, err := svc.ResendEmailConfirmationEmail("")
	require.False(t, sent)
	require.Equal(t, config.ErrEmptyEmail, err)
}

func TestResendEmailConfirmationEmailShouldReturnErrorWhenGetUserIdByEmailFails(t *testing.T) {
	svc, mock := GetServiceAndMock()
	mock.On(GetuserIdByEmail, "asd").Return(int64(0), errors.New("Indistinct"))
	sent, err := svc.ResendEmailConfirmationEmail("asd")
	require.False(t, sent)
	require.Equal(t, "Indistinct", err.Error())
}

func TestResendEmailConfirmationEmailShouldReturnErrorWhenSendConfirmationEmailFaisl(t *testing.T) {
	svc, mock := GetServiceAndMock()
	mock.On(GetuserIdByEmail, "asd").Return(int64(1), nil)
	verified := &UserEmail{
		VerfiedEmail: false,
	}
	mock.On(GetEmailByUserId, int64(1)).Return("notempty", verified, nil)
	sent, err := svc.ResendEmailConfirmationEmail("asd")
	require.False(t, sent)
	require.Equal(t, config.ErrEmailSendServiceNotWorking, err)
}

func TestSendUsernameShouldReturnErrorWhenEmailIsEmpty(t *testing.T) {
	svc, _ := GetServiceAndMock()

	sent, err := svc.SendUsername("")
	require.Equal(t, config.ErrEmptyEmail, err)
	require.False(t, sent)
}

func TestSendUsernameShouldReturnErrorWhenGetUsernameFails(t *testing.T) {
	svc, mock := GetServiceAndMock()
	mock.On(GetUsernameByEmail, "asd").Return("", errors.New("Indistinct"))
	sent, err := svc.SendUsername("asd")
	require.Equal(t, "Indistinct", err.Error())
	require.False(t, sent)
}

func TestSendUsernameShouldReturnErrorWhenEmailSenderFails(t *testing.T) {
	svc, mock := GetServiceAndMock()
	mock.On(GetUsernameByEmail, "asd").Return("juan", nil)
	sent, err := svc.SendUsername("asd")
	require.Equal(t, config.ErrEmailSendServiceNotWorking, err)
	require.False(t, sent)
}