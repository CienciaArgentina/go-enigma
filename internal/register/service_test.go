package register

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type RegisterRepositoryMock struct {
	mock.Mock
}

const (
	AddUser             = "AddUser"
	AddEmail            = "AddEmail"
	VerifyIfEmailExists = "VerifyIfEmailExists"
)

func (r *RegisterRepositoryMock) AddUser(u *User) (int64, error) {
	args := r.Called(u)
	return args.Get(0).(int64), args.Error(1)
}

func (r *RegisterRepositoryMock) AddEmail(e *UserEmail) (int64, error) {
	args := r.Called(e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *RegisterRepositoryMock) VerifyIfEmailExists(email string) (bool, error) {
	args := r.Called(email)
	return args.Get(0).(bool), args.Error(1)
}

func TestDefaultRegisterOptionsShouldReturnOptions(t *testing.T) {
	opt := defaultRegisterOptions()
	require.True(t, opt.UserOptions.RequireUniqueEmail)
}

func TestNewShouldReturnNewService(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	require.NotNil(t, srv)
}

func TestSignUpShouldReturnErrorWhenUserSignUpDtoIsNil(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{}

	errs := srv.SignUp(&u)

	require.Equal(t, errEmptyUsername, errs[0])
}

func TestSignUpShouldReturnErrorWhenUsernameIsEmpty(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "",
	}

	errs := srv.SignUp(&u)

	require.Equal(t, errEmptyUsername, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordIsEmpty(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan",
		Password: "",
	}

	errs := srv.SignUp(&u)

	require.Equal(t, errEmptyPassword, errs[0])
}

func TestSignUpShouldReturnErrorWhenEmailIsEmpty(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan",
		Password: "password",
		Email:    "",
	}

	errs := srv.SignUp(&u)

	require.Equal(t, errEmptyEmail, errs[0])
}

func TestSignUpShouldReturnTrueWhenEmailAlreadyExistsEmpty(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan",
		Password: "password",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(true, nil)

	errs := srv.SignUp(&u)

	require.Equal(t, errEmailAlreadyRegistered, errs[0])
}

func TestSignUpShouldReturnErrorWhenEmailCheckFails(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan",
		Password: "password",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(true, errors.New("Indiferent"))

	errs := srv.SignUp(&u)

	require.Equal(t, "Indiferent", errs[0].Error())
}

func TestSignUpShouldReturnErrorWhenUsingAnInvalidCharInUsername(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan^",
		Password: "password",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	errs := srv.SignUp(&u)

	require.Equal(t, errUsernameCotainsIlegalChars, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordContainsSpace(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan",
		Password: "pass word",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	errs := srv.SignUp(&u)

	require.Equal(t, errPwContainsSpace, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordDoesNotContainsUppercase(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan",
		Password: "password$1",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	errs := srv.SignUp(&u)

	require.Equal(t, errPwDoesNotContainsUppercase, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordDoesNotContainsLowercase(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan",
		Password: "AAAAAAAA$1",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	errs := srv.SignUp(&u)

	require.Equal(t, errPwDoesNotContainsLowercase, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordDoesNotNonAlphanumericChar(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := New(repoMock, nil)
	u := UserSignUpDto{
		Username: "Juan",
		Password: "AAAAAAAAaa1",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	errs := srv.SignUp(&u)

	require.Equal(t, errPwDoesNotContainsNonAlphaChars, errs[0])
}
