package recovery

import (
	"errors"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	GetEmailByUserId = "GetEmailByUserId"
	ConfirmUserEmail = "ConfirmUserEmail"
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
		VerfiedEmail:     true,
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