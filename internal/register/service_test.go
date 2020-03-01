package register

import (
	"errors"
	"github.com/CienciaArgentina/go-enigma/config"
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

var (
	completeUserDto = &UserSignUp{
		Username: "Juancito123",
		Password: "Password!123",
		Email:    "n@n.com",
	}
	cfg = config.DefaultConfiguration()
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
	srv := NewService(repoMock, nil, cfg)
	require.NotNil(t, srv)
}

func TestSignUpShouldReturnErrorWhenUserSignUpDtoIsNil(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{}

	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrEmptyUsername, errs[0])
}

func TestSignUpShouldReturnErrorWhenUsernameIsEmpty(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "",
	}

	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrEmptyUsername, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordIsEmpty(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "",
	}

	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrEmptyPassword, errs[0])
}

func TestSignUpShouldReturnErrorWhenEmailIsEmpty(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "password",
		Email:    "",
	}

	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrEmptyEmail, errs[0])
}

func TestSignUpShouldReturnErrorWhenEmailFormatIsInvalid(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "password",
		Email:    "asd@gm_ail.com",
	}

	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrInvalidEmail, errs[0])
}

func TestSignUpShouldReturnTrueWhenEmailAlreadyExistsEmpty(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "password",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(true, nil)

	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrEmailAlreadyRegistered, errs[0])
}

func TestSignUpShouldReturnErrorWhenEmailCheckFails(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "password",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(true, errors.New("Indiferent"))

	_, errs := srv.SignUp(&u)

	require.Equal(t, "Indiferent", errs[0].Error())
}

func TestSignUpShouldReturnErrorWhenUsingAnInvalidCharInUsername(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan^",
		Password: "password",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrUsernameCotainsIlegalChars, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordContainsSpace(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "pass word",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrPwContainsSpace, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordContainsLessThan8Chars(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "p",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	_, errs := srv.SignUp(&u)

	require.Equal(t, "El campo de contraseña tiene menos de 8 caracteres", errs[0].Error())
}

func TestSignUpShouldReturnErrorWhenPasswordDoesNotContainsUppercase(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "password$1",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	_, errs := srv.SignUp(&u)

	require.Equal(t, config.ErrPwDoesNotContainsUppercase, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordDoesNotContainsLowercase(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "AAAAAAAA$1",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	_, errs := srv.SignUp(&u)

	require.Len(t, errs, 1)
	require.Equal(t, config.ErrPwDoesNotContainsLowercase, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordDoesNotContainsNonAlphanumericChar(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "AAAAAAAAaa1",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	_, errs := srv.SignUp(&u)

	require.Len(t, errs, 1)
	require.Equal(t, config.ErrPwDoesNotContainsNonAlphaChars, errs[0])
}

func TestSignUpShouldReturnErrorWhenPasswordDoesNotContainsADigit(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	u := UserSignUp{
		Username: "Juan",
		Password: "AAaa!!AA",
		Email:    "n@n.com",
	}
	repoMock.On(VerifyIfEmailExists, u.Email).Return(false, nil)
	_, errs := srv.SignUp(&u)

	require.Len(t, errs, 1)
	require.Equal(t, config.ErrPwDoesNotContainsADigit, errs[0])
}

func TestSignUpShouldReturnErrorIfAddUserFails(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	repoMock.On(VerifyIfEmailExists, completeUserDto.Email).Return(false, nil)
	repoMock.On(AddUser, mock.Anything).Return(int64(0), errors.New("Indiferent"))

	_, errs := srv.SignUp(completeUserDto)

	require.Len(t, errs, 1)
	require.Equal(t, "Indiferent", errs[0].Error())
}

func TestSignUpShouldReturnErrorIfAddEmailFails(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	userId := int64(12)
	repoMock.On(VerifyIfEmailExists, completeUserDto.Email).Return(false, nil)
	repoMock.On(AddUser, mock.Anything).Return(userId, nil)
	repoMock.On(AddEmail, mock.Anything).Return(int64(0), errors.New("Indiferent"))

	_, errs := srv.SignUp(completeUserDto)

	require.Len(t, errs, 1)
	require.Equal(t, "Indiferent", errs[0].Error())
}

func TestSignUpShouldNoFail(t *testing.T) {
	repoMock := new(RegisterRepositoryMock)
	srv := NewService(repoMock, nil, cfg)
	userId := int64(12)
	repoMock.On(VerifyIfEmailExists, completeUserDto.Email).Return(false, nil)
	repoMock.On(AddUser, mock.Anything).Return(userId, nil)
	repoMock.On(AddEmail, mock.Anything).Return(int64(0), nil)

	_, errs := srv.SignUp(completeUserDto)

	require.Len(t, errs, 0)
	require.Nil(t, errs)
}