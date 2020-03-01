package login

import (
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type LoginRepositoryMock struct{
	mock.Mock
}

func (l *LoginRepositoryMock) GetUserByUsername(username string) (*User, *UserEmail, error) {
	args := l.Called(username)
	return args.Get(0).(*User), args.Get(1).(*UserEmail), args.Error(2)
}

func (l *LoginRepositoryMock) IncrementLoginFailAttempt(userId int) error {
	args := l.Called(userId)
	return args.Error(0)
}

func (l *LoginRepositoryMock) ResetLoginFails(userId int) error {
	args := l.Called(userId)
	return args.Error(0)
}

func (l *LoginRepositoryMock) UnlockAccount(userId int) error {
	args := l.Called(userId)
	return args.Error(0)
}

func (l *LoginRepositoryMock) LockAccount(userId int, duration time.Duration) error {
	args := l.Called(userId)
	return args.Error(0)
}

func (l *LoginRepositoryMock) GetUserRole(userId int) (string, error) {
	args := l.Called(userId)
	return args.Get(0).(string), args.Error(1)
}

func CreateNewServiceWithMockedRepo() Service {
	repo := new(LoginRepositoryMock)
	cfg := config.New()
	return NewService(repo, nil, cfg)
}

func TestNewServiceShouldReturnNewService(t *testing.T) {
	svc := CreateNewServiceWithMockedRepo()
	require.NotNil(t, svc)
}

func TestDefaultLoginOptionsShouldReturnDefaultOptions(t *testing.T) {
	opt := defaultLoginOptions()
	require.NotNil(t, opt.LockoutOptions.LockoutTimeDuration)
}

func TestVerifyCanLoginShouldFailIfUsernameIsEmpty(t *testing.T) {
	svc := CreateNewServiceWithMockedRepo()
	var userLogin UserLogin
	login, err := svc.VerifyCanLogin(&userLogin)
	require.Equal(t, config.ErrEmptyUsername, err)
	require.False(t, login)
}

func TestVerifyCanLoginShouldFailIfPasswordIsEmpty(t *testing.T) {
	svc := CreateNewServiceWithMockedRepo()
	var userLogin UserLogin
	userLogin.Username = "notempty"
	login, err := svc.VerifyCanLogin(&userLogin)
	require.Equal(t, config.ErrEmptyPassword, err)
	require.False(t, login)
}

func TestLoginShouldFailIfUsernameOrPasswordAreEmpty(t *testing.T) {
	svc := CreateNewServiceWithMockedRepo()
	var userLogin UserLogin

	login, err := svc.Login(&userLogin)
	require.Equal(t, config.ErrEmptyUsername, err)
	require.Equal(t, "",login)

	userLogin.Username = "notempty"
	login, err = svc.Login(&userLogin)
	require.Equal(t, config.ErrEmptyPassword, err)
	require.Equal(t, "",login)
}