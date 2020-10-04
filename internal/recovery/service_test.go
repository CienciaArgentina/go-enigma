package recovery

import (
	"errors"
	"reflect"
	"testing"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	domain2 "github.com/CienciaArgentina/go-enigma/internal/domain"
)

const (
	GetEmailByUserIdMockID = iota
	ConfirmUserEmailMockID
	GetuserIdByEmailMockID
	GetUsernameByEmailMockID
	GetSecurityTokenMockID
	UpdatePasswordHashMockID
	UpdateSecurityTokenMockID
)

type MockRepository struct {
	Responses map[int]interface{}
	Errors    map[int]apierror.ApiError
}

func (m *MockRepository) GetEmailByUserId(userId int64) (string, *domain2.UserEmail, apierror.ApiError) {
	return "token", m.Responses[GetEmailByUserIdMockID].(*domain2.UserEmail), m.Errors[GetEmailByUserIdMockID]
}

func (m *MockRepository) ConfirmUserEmail(email string, token string) apierror.ApiError {
	return m.Errors[ConfirmUserEmailMockID]
}

func (m *MockRepository) GetuserIdByEmail(email string) (int64, apierror.ApiError) {
	return m.Responses[GetuserIdByEmailMockID].(int64), m.Errors[GetuserIdByEmailMockID]
}

func (m *MockRepository) GetUsernameByEmail(email string) (string, apierror.ApiError) {
	return m.Responses[GetUsernameByEmailMockID].(string), m.Errors[GetUsernameByEmailMockID]
}

func (m *MockRepository) GetSecurityToken(email string) (string, apierror.ApiError) {
	return m.Responses[GetSecurityTokenMockID].(string), m.Errors[GetSecurityTokenMockID]
}

func (m *MockRepository) UpdatePasswordHash(userId int64, passwordHash string) (bool, apierror.ApiError) {
	return m.Responses[UpdatePasswordHashMockID].(bool), m.Errors[UpdatePasswordHashMockID]
}

func (m *MockRepository) UpdateSecurityToken(userId int64, newSecurityToken string) (bool, apierror.ApiError) {
	return m.Responses[UpdateSecurityTokenMockID].(bool), m.Errors[UpdateSecurityTokenMockID]
}

func (m *MockRepository) GetUserByUserId(userId int64) (*domain2.User, apierror.ApiError) {
	return m.Responses[GetUserByUserIdMockID].(*domain2.User), m.Errors[GetUserByUserIdMockID]
}

func Test_recoveryService_GetUserByUserId(t *testing.T) {
	type fields struct {
		repository RecoveryRepository
		cfg        *config.EnigmaConfig
	}

	type args struct {
		userId int64
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.User
		want1  apierror.ApiError
	}{
		{
			name: "ok",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetUserByUserIdMockID: &domain2.User{Username: "koppin"},
					},
					Errors: map[int]apierror.ApiError{
						GetUserByUserIdMockID: nil,
					},
				},
			},
			args:  args{},
			want:  &domain2.User{Username: "koppin"},
			want1: nil,
		},
		{
			name: "error",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetUserByUserIdMockID: &domain2.User{Username: "koppin"},
					},
					Errors: map[int]apierror.ApiError{
						GetUserByUserIdMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			args:  args{},
			want:  nil,
			want1: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &recoveryService{
				repository: tt.fields.repository,
				cfg:        tt.fields.cfg,
			}
			got, got1 := r.GetUserByUserId(tt.args.userId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("recoveryService.GetUserByUserId() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("recoveryService.GetUserByUserId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_recoveryService_ConfirmEmail(t *testing.T) {
	type fields struct {
		repository RecoveryRepository
		cfg        *config.EnigmaConfig
	}

	type args struct {
		email string
		token string
		ctx   *rest.ContextInformation
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  apierror.ApiError
	}{
		{
			name: "empty_email",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						ConfirmUserEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmUserEmailMockID: nil,
					},
				},
			},
			args: args{
				token: "test",
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(ErrEmailValidationFailed),
		},
		{
			name: "empty_token",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						ConfirmUserEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmUserEmailMockID: nil,
					},
				},
			},
			args: args{
				email: "test",
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(ErrEmailValidationFailed),
		},
		{
			name: "empty_email_token",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						ConfirmUserEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmUserEmailMockID: nil,
					},
				},
			},
			args:  args{},
			want:  false,
			want1: apierror.NewBadRequestApiError(ErrEmailValidationFailed),
		},
		{
			name: "internal_error",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						ConfirmUserEmailMockID: false,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmUserEmailMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			args: args{
				email: "test",
				token: "test",
				ctx:   &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
		},
		{
			name: "ok",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						ConfirmUserEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmUserEmailMockID: nil,
					},
				},
			},
			args: args{
				email: "test",
				token: "test",
				ctx:   &rest.ContextInformation{},
			},
			want:  true,
			want1: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &recoveryService{
				repository: tt.fields.repository,
				cfg:        tt.fields.cfg,
			}
			got, got1 := r.ConfirmEmail(tt.args.email, tt.args.token, tt.args.ctx)
			if got != tt.want {
				t.Errorf("recoveryService.ConfirmEmail() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("recoveryService.ConfirmEmail() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_recoveryService_ResendEmailConfirmationEmail(t *testing.T) {
	type fields struct {
		repository RecoveryRepository
		cfg        *config.EnigmaConfig
	}
	type args struct {
		email string
		ctx   *rest.ContextInformation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  apierror.ApiError
	}{
		{
			name: "empty_email",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetEmailByUserIdMockID: &domain.UserEmail{
							UserId: 123,
						},
						GetuserIdByEmailMockID: int64(123),
					},
					Errors: map[int]apierror.ApiError{
						GetEmailByUserIdMockID: nil,
					},
				},
			},
			args: args{
				email: "",
				ctx:   &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(domain.ErrEmptyEmail),
		},
		{
			name: "get_user_error",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetEmailByUserIdMockID: &domain.UserEmail{
							UserId: 123,
						},
						GetuserIdByEmailMockID: int64(123),
					},
					Errors: map[int]apierror.ApiError{
						GetEmailByUserIdMockID: nil,
						GetuserIdByEmailMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			args: args{
				email: "test",
				ctx:   &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
		},
		{
			name: "ok",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetEmailByUserIdMockID: &domain.UserEmail{
							UserId: 123,
						},
						GetuserIdByEmailMockID: int64(123),
					},
					Errors: map[int]apierror.ApiError{
						GetEmailByUserIdMockID: nil,
					},
				},
			},
			args: args{
				email: "test",
				ctx:   &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError("cant send email", errors.New("cant send email"), "cannot_email"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &recoveryService{
				repository: tt.fields.repository,
				cfg:        tt.fields.cfg,
			}
			got, got1 := r.ResendEmailConfirmationEmail(tt.args.email, tt.args.ctx)
			if got != tt.want {
				t.Errorf("recoveryService.ResendEmailConfirmationEmail() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("recoveryService.ResendEmailConfirmationEmail() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_recoveryService_SendConfirmationEmail(t *testing.T) {
	type fields struct {
		repository RecoveryRepository
		cfg        *config.EnigmaConfig
	}
	type args struct {
		userId int64
		ctx    *rest.ContextInformation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  apierror.ApiError
	}{
		{
			name: "internal_error_getting_user",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetEmailByUserIdMockID: &domain2.UserEmail{
							UserId: 123,
						},
					},
					Errors: map[int]apierror.ApiError{
						GetEmailByUserIdMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			args: args{
				userId: 123,
				ctx:    &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
		},
		{
			name: "empty_user_email",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetEmailByUserIdMockID: &domain.UserEmail{},
					},
					Errors: map[int]apierror.ApiError{
						GetEmailByUserIdMockID: nil,
					},
				},
			},
			args: args{
				userId: 123,
				ctx:    &rest.ContextInformation{},
			},
			want:  true,
			want1: nil,
		},
		{
			name: "verified_email",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetEmailByUserIdMockID: &domain.UserEmail{
							VerfiedEmail: true,
						},
					},
					Errors: map[int]apierror.ApiError{
						GetEmailByUserIdMockID: nil,
					},
				},
			},
			args: args{
				userId: 123,
				ctx:    &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(ErrEmailAlreadyVerified),
		},
		{
			name: "ok",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetEmailByUserIdMockID: &domain.UserEmail{
							UserId: 123,
						},
					},
					Errors: map[int]apierror.ApiError{
						GetEmailByUserIdMockID: nil,
					},
				},
			},
			args: args{
				userId: 123,
				ctx:    &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError("cant send email", errors.New("cant send email"), "cannot_email"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &recoveryService{
				repository: tt.fields.repository,
				cfg:        tt.fields.cfg,
			}
			got, got1 := r.SendConfirmationEmail(tt.args.userId, tt.args.ctx)
			if got != tt.want {
				t.Errorf("recoveryService.SendConfirmationEmail() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("recoveryService.SendConfirmationEmail() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_recoveryService_SendUsername(t *testing.T) {
	type fields struct {
		repository RecoveryRepository
		cfg        *config.EnigmaConfig
	}
	type args struct {
		email string
		ctx   *rest.ContextInformation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  apierror.ApiError
	}{
		{
			name: "empty_email",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetUsernameByEmailMockID: &domain.UserEmail{
							UserId: 123,
						},
					},
					Errors: map[int]apierror.ApiError{},
				},
			},
			args: args{
				email: "",
				ctx:   &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewBadRequestApiError(domain.ErrEmptyEmail),
		},
		{
			name: "error_username",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetUsernameByEmailMockID: "test",
					},
					Errors: map[int]apierror.ApiError{
						GetUsernameByEmailMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			args: args{
				email: "test",
				ctx:   &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
		},
		{
			name: "ok",
			fields: fields{
				repository: &MockRepository{
					Responses: map[int]interface{}{
						GetUsernameByEmailMockID: "test",
					},
					Errors: map[int]apierror.ApiError{
						GetUsernameByEmailMockID: nil,
					},
				},
			},
			args: args{
				email: "test",
				ctx:   &rest.ContextInformation{},
			},
			want:  false,
			want1: apierror.NewInternalServerApiError("cant send email", errors.New("cant send email"), "cannot_email"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &recoveryService{
				repository: tt.fields.repository,
				cfg:        tt.fields.cfg,
			}
			got, got1 := r.SendUsername(tt.args.email, tt.args.ctx)
			if got != tt.want {
				t.Errorf("recoveryService.SendUsername() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("recoveryService.SendUsername() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
