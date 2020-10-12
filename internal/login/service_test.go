package login

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
)

const (
	GetUserByUsernameMockID = iota
	IncrementLoginFailAttemptMockID
	ResetLoginFailsMockID
	UnlockAccountMockID
	LockAccountMockID
)

type MockRepository struct {
	Responses map[int]interface{}
	Errors    map[int]apierror.ApiError
}

func (m *MockRepository) GetUserByUsername(username string) (*domain.User, *domain.UserEmail, apierror.ApiError) {
	return m.Responses[GetUserByUsernameMockID].([]interface{})[0].(*domain.User),
		m.Responses[GetUserByUsernameMockID].([]interface{})[1].(*domain.UserEmail),
		m.Errors[GetUserByUsernameMockID]
}

func (m *MockRepository) IncrementLoginFailAttempt(userID int64) error {
	return m.Errors[IncrementLoginFailAttemptMockID]
}

func (m *MockRepository) ResetLoginFails(userID int64) error {
	return m.Errors[ResetLoginFailsMockID]
}

func (m *MockRepository) UnlockAccount(userID int64) error {
	return m.Errors[UnlockAccountMockID]
}

func (m *MockRepository) LockAccount(userID int64, duration time.Duration) error {
	return m.Errors[LockAccountMockID]
}

func Test_loginService_LoginUser(t *testing.T) {
	type fields struct {
		cfg          *config.EnigmaConfig
		loginOptions *config.LoginOptions
		repository   Repository
	}
	type args struct {
		u   *domain.UserLoginDTO
		ctx *rest.ContextInformation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  apierror.ApiError
	}{
		{
			name: "empty_username",
			args: args{
				u: &domain.UserLoginDTO{
					Username: "",
					Password: "test",
				},
				ctx: &rest.ContextInformation{},
			},
			want:  "",
			want1: apierror.NewBadRequestApiError(domain.ErrEmptyUsername),
		},
		{
			name: "empty_password",
			args: args{
				u: &domain.UserLoginDTO{
					Username: "test",
					Password: "",
				},
				ctx: &rest.ContextInformation{},
			},
			want:  "",
			want1: apierror.NewBadRequestApiError(domain.ErrEmptyPassword),
		},
		{
			name: "error_getting_user",
			args: args{
				u: &domain.UserLoginDTO{
					Username: "test",
					Password: "test",
				},
				ctx: &rest.ContextInformation{},
			},
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetUserByUsernameMockID: []interface{}{
							&domain.User{},
							&domain.UserEmail{},
						},
					},
					Errors: map[int]apierror.ApiError{
						GetUserByUsernameMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			want:  "",
			want1: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
		},
		{
			name: "password_mismatch",
			args: args{
				u: &domain.UserLoginDTO{
					Username: "test",
					Password: "test",
				},
				ctx: &rest.ContextInformation{},
			},
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetUserByUsernameMockID: []interface{}{
							&domain.User{
								PasswordHash: "aisjdoajsid",
							},
							&domain.UserEmail{},
						},
					},
					Errors: map[int]apierror.ApiError{
						GetUserByUsernameMockID: nil,
					},
				},
			},
			want:  "",
			want1: apierror.NewInternalServerApiError(domain.ErrUnexpectedError, errors.New("encoded hash string is not 6"), domain.ErrInternalCode),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &loginService{
				cfg:          tt.fields.cfg,
				loginOptions: tt.fields.loginOptions,
				repository:   tt.fields.repository,
			}
			got, got1 := l.LoginUser(tt.args.u, tt.args.ctx)
			if got != tt.want {
				t.Errorf("loginService.LoginUser() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("loginService.LoginUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
