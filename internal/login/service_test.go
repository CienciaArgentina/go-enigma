package login

import (
	"database/sql"
	"errors"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	GetUserByUsername         = "GetUserByUsername"
	IncrementLoginFailAttempt = "IncrementLoginFailAttempt"
	ResetLoginFails           = "ResetLoginFails"
	UnlockAccount             = "UnlockAccount"
	LockAccount               = "LockAccount"
	GetUserRole               = "GetUserRole"
)

type LoginRepositoryMock struct {
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

func CreateNewServiceWithMockedRepo() (Service, *LoginRepositoryMock) {
	repo := new(LoginRepositoryMock)
	cfg := config.New()
	return NewService(repo, nil, cfg), repo
}

func ReturnFullUserDto() *UserLogin {
	return &UserLogin{
		Username: "Juan",
		Password: "Hola!123*",
	}
}

func ReturnFullUser() *User {
	return &User{
		UserId:              1,
		Username:            "Juancito",
		NormalizedUsername:  "Juancito",
		PasswordHash:        "asdf",
		LockoutEnabled:      false,
		LockoutDate:         mysql.NullTime{},
		FailedLoginAttempts: 0,
		DateCreated:         "",
		SecurityToken:       sql.NullString{},
		VerificationToken:   "",
		DateDeleted:         nil,
	}
}

func ReturnFullUserEmail() *UserEmail {
	return &UserEmail{
		UserEmailId:      1,
		UserId:           1,
		Email:            "asd@asd.com",
		NormalizedEmail:  "ASD@ASD.COM",
		VerfiedEmail:     false,
		VerificationDate: nil,
		DateCreated:      "",
		DateDeleted:      sql.NullTime{},
	}
}

func TestNewServiceShouldReturnNewService(t *testing.T) {
	svc, _ := CreateNewServiceWithMockedRepo()
	require.NotNil(t, svc)
}

func TestDefaultLoginOptionsShouldReturnDefaultOptions(t *testing.T) {
	opt := defaultLoginOptions()
	require.NotNil(t, opt.LockoutOptions.LockoutTimeDuration)
}

func TestVerifyCanLoginShouldFailIfUsernameIsEmpty(t *testing.T) {
	svc, _ := CreateNewServiceWithMockedRepo()
	var userLogin UserLogin
	login, err := svc.VerifyCanLogin(&userLogin)
	require.Equal(t, config.ErrEmptyUsername, err)
	require.False(t, login)
}

func TestVerifyCanLoginShouldFailIfPasswordIsEmpty(t *testing.T) {
	svc, _ := CreateNewServiceWithMockedRepo()
	var userLogin UserLogin
	userLogin.Username = "notempty"
	login, err := svc.VerifyCanLogin(&userLogin)
	require.Equal(t, config.ErrEmptyPassword, err)
	require.False(t, login)
}

func TestLoginShouldFailIfUsernameOrPasswordAreEmpty(t *testing.T) {
	svc, _ := CreateNewServiceWithMockedRepo()
	var userLogin UserLogin

	login, err := svc.Login(&userLogin)
	require.Equal(t, config.ErrEmptyUsername, err)
	require.Equal(t, "", login)

	userLogin.Username = "notempty"
	login, err = svc.Login(&userLogin)
	require.Equal(t, config.ErrEmptyPassword, err)
	require.Equal(t, "", login)
}

func TestLoginShouldThrowInvalidLoginIfTheresAndErrorGettingUsername(t *testing.T) {
	svc, mock := CreateNewServiceWithMockedRepo()
	mock.On(GetUserByUsername, ReturnFullUserDto().Username).Return(&User{}, &UserEmail{}, errors.New("Indistinct"))

	login, err := svc.Login(ReturnFullUserDto())

	require.Error(t, err)
	require.Empty(t, login)
}

func TestLoginShouldThrowInvalidLoginIfUserIsEmpty(t *testing.T) {
	svc, mock := CreateNewServiceWithMockedRepo()
	mock.On(GetUserByUsername, ReturnFullUserDto().Username).Return(&User{}, &UserEmail{}, nil)

	login, err := svc.Login(ReturnFullUserDto())

	require.Error(t, err)
	require.Empty(t, login)
}

func TestLoginShouldThrowErrorWhileComparingNonIdenticalPasswords(t *testing.T) {
	svc, mock := CreateNewServiceWithMockedRepo()
	mock.On(GetUserByUsername, ReturnFullUserDto().Username).Return(ReturnFullUser(), ReturnFullUserEmail(), nil)

	login, err := svc.Login(ReturnFullUserDto())

	require.Equal(t, config.ErrThroughLogin, err)
	require.Empty(t, login)
}
