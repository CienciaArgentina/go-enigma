package recovery

import (
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type RecoveryRepositoryMock struct {
	mock.Mock
}

func (r *RecoveryRepositoryMock) GetEmailByUserId(userId int64) (string, *UserEmail, error) {
	args := r.Called(userId)
	return args.Get(0).(string), args.Get(1).(*UserEmail), args.Error(2)
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