package register

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	GetUserByIDMockName         = "GetUserById"
	AddUserMockName             = "AddUser"
	AddUserEmailMockName        = "AddUserEmail"
	DeleteUserMockName          = "DeleteUser"
	CheckUsernameExistsMockName = "CheckUsernameExists"
	CheckEmailExistsMockName    = "CheckEmailExists"
)

type MockRepository struct {
	Responses map[string]interface{}
	Errors    map[string]error
}

func (m *MockRepository) GetUserById(userId int64) (*domain.User, error) {
	return m.Responses[GetUserByIDMockName].(*domain.User), m.Errors[GetUserByIDMockName]
}

func (m *MockRepository) AddUser(tx *sqlx.Tx, u *domain.User) (int64, error) {
	return m.Responses[AddUserMockName].(int64), m.Errors[AddUserMockName]
}

func (m *MockRepository) AddUserEmail(tx *sqlx.Tx, e *domain.UserEmail) (int64, error) {
	return m.Responses[AddUserEmailMockName].(int64), m.Errors[AddUserEmailMockName]
}

func (m *MockRepository) DeleteUser(userId int64) error {
	return m.Errors[DeleteUserMockName]
}

func (m *MockRepository) CheckUsernameExists(username string) (bool, error) {
	return m.Responses[CheckUsernameExistsMockName].(bool), m.Errors[CheckUsernameExistsMockName]
}

func (m *MockRepository) CheckEmailExists(email string) (bool, error) {
	return m.Responses[CheckEmailExistsMockName].(bool), m.Errors[CheckEmailExistsMockName]
}

func Test_registerService_UserCanSignUp(t *testing.T) {
	type fields struct {
		repository RegisterRepository
	}

	type args struct {
		usr *domain.UserSignupDTO
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  apierror.ApiError
	}{
		{
			name: "username_empty",
			args: args{
				usr: &domain.UserSignupDTO{
					Password: "test",
					Email:    "test",
				},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(domain.ErrEmptyUsername),
		},
		{
			name: "password_empty",
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test",
					Email:    "test",
				},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(domain.ErrEmptyPassword),
		},
		{
			name: "email_empty",
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test",
					Password: "test",
				},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(domain.ErrEmptyEmail),
		},
		{
			name: "invalid_email",
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test",
					Password: "test",
					Email:    "test",
				},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(errInvalidEmailFormat),
		},
		{
			name: "non_unique_email",
			fields: fields{
				repository: &MockRepository{
					Responses: map[string]interface{}{
						CheckEmailExistsMockName: true,
					},
				},
			},
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test",
					Password: "test",
					Email:    "test@gmail.com",
				},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(errEmailAlreadyExists),
		},
		{
			name: "check_email_exists_error",
			fields: fields{
				repository: &MockRepository{
					Responses: map[string]interface{}{
						CheckEmailExistsMockName: false,
					},
					Errors: map[string]error{
						CheckEmailExistsMockName: errCannotDelete,
					},
				},
			},
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test",
					Password: "test",
					Email:    "test@gmail.com",
				},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError(errEmailAlreadyExistsInternalErr, errCannotDelete, domain.ErrInternalCode),
		},
		{
			name: "check_user_exists",
			fields: fields{
				repository: &MockRepository{
					Responses: map[string]interface{}{
						CheckEmailExistsMockName:    false,
						CheckUsernameExistsMockName: true,
					},
					Errors: map[string]error{
						CheckEmailExistsMockName: sql.ErrNoRows,
					},
				},
			},
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test",
					Password: "test",
					Email:    "test@gmail.com",
				},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(errUserAlreadyExists),
		},
		{
			name: "check_user_exists_error",
			fields: fields{
				repository: &MockRepository{
					Responses: map[string]interface{}{
						CheckEmailExistsMockName:    false,
						CheckUsernameExistsMockName: false,
					},
					Errors: map[string]error{
						CheckEmailExistsMockName:    sql.ErrNoRows,
						CheckUsernameExistsMockName: errCannotDelete,
					},
				},
			},
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test",
					Password: "test",
					Email:    "test@gmail.com",
				},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError(errUserAlreadyExistsInternalErr, errCannotDelete, domain.ErrInternalCode),
		},
		{
			name: "invalid_user_request_lower",
			fields: fields{
				repository: &MockRepository{
					Responses: map[string]interface{}{
						CheckEmailExistsMockName:    false,
						CheckUsernameExistsMockName: false,
					},
					Errors: map[string]error{
						CheckEmailExistsMockName:    sql.ErrNoRows,
						CheckUsernameExistsMockName: sql.ErrNoRows,
					},
				},
			},
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test?",
					Password: "test ",
					Email:    "test@gmail.com",
				},
			},
			want: false,
			want1: apierror.
				NewWithStatus(http.StatusBadRequest).
				WithMessage(errCantCreateUser).
				AddError(errUsernameCotainsIlegalChars, errInvalidUsernameCode).
				AddError(errPwContainsSpace, errInvalidPasswordCode).
				AddError(fmt.Sprintf("El campo de contraseña tiene menos de %d caracteres", 8), errInvalidPasswordCode).
				AddError(errPwDoesNotContainsUppercase, errInvalidPasswordCode).
				AddError(errPwDoesNotContainsNonAlphaChars, errInvalidPasswordCode).
				AddError(errPwDoesNotContainsADigit, errInvalidPasswordCode),
		},
		{
			name: "invalid_user_request_lower",
			fields: fields{
				repository: &MockRepository{
					Responses: map[string]interface{}{
						CheckEmailExistsMockName:    false,
						CheckUsernameExistsMockName: false,
					},
					Errors: map[string]error{
						CheckEmailExistsMockName:    sql.ErrNoRows,
						CheckUsernameExistsMockName: sql.ErrNoRows,
					},
				},
			},
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test?",
					Password: "TEST ",
					Email:    "test@gmail.com",
				},
			},
			want: false,
			want1: apierror.
				NewWithStatus(http.StatusBadRequest).
				WithMessage(errCantCreateUser).
				AddError(errUsernameCotainsIlegalChars, errInvalidUsernameCode).
				AddError(errPwContainsSpace, errInvalidPasswordCode).
				AddError(fmt.Sprintf("El campo de contraseña tiene menos de %d caracteres", 8), errInvalidPasswordCode).
				AddError(errPwDoesNotContainsLowercase, errInvalidPasswordCode).
				AddError(errPwDoesNotContainsNonAlphaChars, errInvalidPasswordCode).
				AddError(errPwDoesNotContainsADigit, errInvalidPasswordCode),
		},
		{
			name: "invalid_user_request_lower",
			fields: fields{
				repository: &MockRepository{
					Responses: map[string]interface{}{
						CheckEmailExistsMockName:    false,
						CheckUsernameExistsMockName: false,
					},
					Errors: map[string]error{
						CheckEmailExistsMockName:    sql.ErrNoRows,
						CheckUsernameExistsMockName: sql.ErrNoRows,
					},
				},
			},
			args: args{
				usr: &domain.UserSignupDTO{
					Username: "test",
					Password: "ThisIsATest123.",
					Email:    "test@gmail.com",
				},
			},
			want:  true,
			want1: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &registerService{
				registerOptions: initRegisterOptions(),
				repository:      tt.fields.repository,
			}
			got, got1 := u.UserCanSignUp(tt.args.usr)
			if got != tt.want {
				t.Errorf("registerService.UserCanSignUp() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("registerService.UserCanSignUp() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
