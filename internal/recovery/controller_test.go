package recovery

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/gin-gonic/gin"
)

const (
	SendConfirmationEmailMockID = iota
	ConfirmEmailMockID
	ResendEmailConfirmationEmailMockID
	SendUsernameMockID
	SendPasswordResetMockID
	ResetPasswordMockID
	GetUserByUserIdMockID
)

// MockService Mock service
type MockService struct {
	Responses map[int]interface{}
	Errors    map[int]apierror.ApiError
}

func (m *MockService) SendConfirmationEmail(userID int64, ctx *rest.ContextInformation) (bool, apierror.ApiError) {
	return m.Responses[SendConfirmationEmailMockID].(bool), m.Errors[SendConfirmationEmailMockID]
}

func (m *MockService) ConfirmEmail(email string, token string, ctx *rest.ContextInformation) (bool, apierror.ApiError) {
	return m.Responses[ConfirmEmailMockID].(bool), m.Errors[ConfirmEmailMockID]
}

func (m *MockService) ResendEmailConfirmationEmail(email string, ctx *rest.ContextInformation) (bool, apierror.ApiError) {
	return m.Responses[ResendEmailConfirmationEmailMockID].(bool), m.Errors[ResendEmailConfirmationEmailMockID]
}

func (m *MockService) SendUsername(email string, ctx *rest.ContextInformation) (bool, apierror.ApiError) {
	return m.Responses[SendUsernameMockID].(bool), m.Errors[SendUsernameMockID]
}

func (m *MockService) SendPasswordReset(email string, ctx *rest.ContextInformation) (bool, apierror.ApiError) {
	return m.Responses[SendPasswordResetMockID].(bool), m.Errors[SendPasswordResetMockID]
}

func (m *MockService) ResetPassword(email, password, confirmPassword, token string, ctx *rest.ContextInformation) (bool, apierror.ApiError) {
	return m.Responses[ResetPasswordMockID].(bool), m.Errors[ResetPasswordMockID]
}

func (m *MockService) GetUserByUserId(userID int64) (*domain.User, apierror.ApiError) {
	return m.Responses[GetUserByUserIdMockID].(*domain.User), m.Errors[GetUserByUserIdMockID]
}

func Test_recoveryController_GetUserByUserId(t *testing.T) {
	type fields struct {
		svc RecoveryService
	}

	type args struct {
		c *gin.Context
	}

	tests := []struct {
		name           string
		fields         fields
		expectedBody   interface{}
		expectedStatus int
		requestBody    interface{}
		params         []gin.Param
	}{
		{
			name: "missing_id",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{},
					Errors:    map[int]apierror.ApiError{},
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   apierror.NewBadRequestApiError(domain.ErrEmptyField),
			params: []gin.Param{
				{
					Key:   "",
					Value: "",
				},
			},
		},
		{
			name: "wrong_id",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{},
					Errors:    map[int]apierror.ApiError{},
				},
			},
			expectedBody:   apierror.NewBadRequestApiError(domain.ErrEmptyField),
			expectedStatus: http.StatusBadRequest,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "test",
				},
			},
		},
		{
			name: "error_not_found",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						GetUserByUserIdMockID: &domain.User{},
					},
					Errors: map[int]apierror.ApiError{
						GetUserByUserIdMockID: apierror.NewNotFoundApiError("User not found"),
					},
				},
			},
			expectedBody:   apierror.NewNotFoundApiError("User not found"),
			expectedStatus: http.StatusNotFound,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "123",
				},
			},
		},
		{
			name: "error_internal",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						GetUserByUserIdMockID: &domain.User{},
					},
					Errors: map[int]apierror.ApiError{
						GetUserByUserIdMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			expectedBody:   apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
			expectedStatus: http.StatusInternalServerError,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "123",
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						GetUserByUserIdMockID: &domain.User{
							Username: "koppin",
						},
					},
					Errors: map[int]apierror.ApiError{
						GetUserByUserIdMockID: nil,
					},
				},
			},
			expectedBody: &domain.User{
				Username: "koppin",
			},
			expectedStatus: http.StatusOK,
			params: []gin.Param{
				{
					Key:   "id",
					Value: "123",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			for _, param := range tt.params {
				c.Params = append(c.Params, param)
			}

			if tt.requestBody != nil {
				req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(tt.requestBody.(string)))
				if err != nil {
					panic(err)
				}
				c.Request = req
			}

			ctr := NewController(tt.fields.svc)

			ctr.GetUserByUserId(c)

			response := w.Result()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("[SignUp] Expected status code = %v, got %v", tt.expectedStatus, response.StatusCode)
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			body := buf.String()

			expected, _ := json.Marshal(tt.expectedBody)
			if !reflect.DeepEqual(body, string(expected)) {
				t.Errorf("[SignUp] Expected body = %v, got %v", tt.expectedBody, body)
				return
			}
		})
	}
}

func Test_recoveryController_ConfirmPasswordReset(t *testing.T) {
	type fields struct {
		svc RecoveryService
	}

	type args struct {
		c *gin.Context
	}

	tests := []struct {
		name           string
		fields         fields
		expectedBody   interface{}
		expectedStatus int
		requestBody    interface{}
		params         []gin.Param
	}{
		{
			name: "bad_request",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{},
					Errors:    map[int]apierror.ApiError{},
				},
			},
			expectedBody:   apierror.NewBadRequestApiError("invalid character 'Â' looking for beginning of value"),
			requestBody:    "´''",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty_body",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{},
					Errors:    map[int]apierror.ApiError{},
				},
			},
			expectedBody:   apierror.NewBadRequestApiError(domain.ErrEmptyField),
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal_server_error",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ResetPasswordMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ResetPasswordMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			expectedBody:   apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
			requestBody:    `{"password": "test"}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ResetPasswordMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ResetPasswordMockID: nil,
					},
				},
			},
			requestBody:    `{"password": "test"}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			for _, param := range tt.params {
				c.Params = append(c.Params, param)
			}

			if tt.requestBody != nil {
				req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(tt.requestBody.(string)))
				if err != nil {
					panic(err)
				}
				c.Request = req
			}

			ctr := NewController(tt.fields.svc)

			ctr.ConfirmPasswordReset(c)

			response := w.Result()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("[SignUp] Expected status code = %v, got %v", tt.expectedStatus, response.StatusCode)
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			body := buf.String()

			expected, _ := json.Marshal(tt.expectedBody)
			if !reflect.DeepEqual(body, string(expected)) && tt.expectedBody != nil {
				t.Errorf("[SignUp] Expected body = %v, got %v", tt.expectedBody, body)
				return
			}
		})
	}
}

func Test_recoveryController_SendPasswordReset(t *testing.T) {
	type fields struct {
		svc RecoveryService
	}

	type args struct {
		c *gin.Context
	}

	tests := []struct {
		name           string
		fields         fields
		expectedBody   interface{}
		expectedStatus int
		URL            string
		requestBody    interface{}
		params         []gin.Param
	}{
		{
			name: "empty_email",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendPasswordResetMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendPasswordResetMockID: nil,
					},
				},
			},
			URL:            "",
			expectedBody:   apierror.NewBadRequestApiError(domain.ErrEmptyField),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal_error",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendPasswordResetMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendPasswordResetMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			URL:            "/test?email=test",
			expectedBody:   apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendPasswordResetMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendPasswordResetMockID: nil,
					},
				},
			},
			URL:            "/test?email=test",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			for _, param := range tt.params {
				c.Params = append(c.Params, param)
			}

			req, err := http.NewRequest(http.MethodGet, tt.URL, nil)
			if tt.requestBody != nil {
				req, err = http.NewRequest(http.MethodPost, tt.URL, strings.NewReader(tt.requestBody.(string)))
			}
			if err != nil {
				panic(err)
			}
			c.Request = req

			ctr := NewController(tt.fields.svc)

			ctr.SendPasswordReset(c)

			response := w.Result()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("[SignUp] Expected status code = %v, got %v", tt.expectedStatus, response.StatusCode)
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			body := buf.String()

			expected, _ := json.Marshal(tt.expectedBody)
			if !reflect.DeepEqual(body, string(expected)) && tt.expectedBody != nil {
				t.Errorf("[SignUp] Expected body = %v, got %v", tt.expectedBody, body)
				return
			}
		})
	}
}

func Test_recoveryController_ForgotUsername(t *testing.T) {
	type fields struct {
		svc RecoveryService
	}

	type args struct {
		c *gin.Context
	}

	tests := []struct {
		name           string
		fields         fields
		expectedBody   interface{}
		expectedStatus int
		URL            string
		requestBody    interface{}
		params         []gin.Param
	}{
		{
			name: "empty_email",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendUsernameMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendUsernameMockID: nil,
					},
				},
			},
			URL:            "",
			expectedBody:   apierror.NewBadRequestApiError(domain.ErrEmptyField),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal_error",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendUsernameMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendUsernameMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			URL:            "/test?email=test",
			expectedBody:   apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendUsernameMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendUsernameMockID: nil,
					},
				},
			},
			URL:            "/test?email=test",
			expectedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			for _, param := range tt.params {
				c.Params = append(c.Params, param)
			}

			req, err := http.NewRequest(http.MethodGet, tt.URL, nil)
			if tt.requestBody != nil {
				req, err = http.NewRequest(http.MethodPost, tt.URL, strings.NewReader(tt.requestBody.(string)))
			}
			if err != nil {
				panic(err)
			}
			c.Request = req

			ctr := NewController(tt.fields.svc)

			ctr.ForgotUsername(c)

			response := w.Result()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("[SignUp] Expected status code = %v, got %v", tt.expectedStatus, response.StatusCode)
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			body := buf.String()

			expected, _ := json.Marshal(tt.expectedBody)
			if !reflect.DeepEqual(body, string(expected)) && tt.expectedBody != nil {
				t.Errorf("[SignUp] Expected body = %v, got %v", tt.expectedBody, body)
				return
			}
		})
	}
}

func Test_recoveryController_ResendEmailConfirmation(t *testing.T) {
	type fields struct {
		svc RecoveryService
	}

	type args struct {
		c *gin.Context
	}

	tests := []struct {
		name           string
		fields         fields
		expectedBody   interface{}
		expectedStatus int
		URL            string
		requestBody    interface{}
		params         []gin.Param
	}{
		{
			name: "empty_email",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ResendEmailConfirmationEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ResendEmailConfirmationEmailMockID: nil,
					},
				},
			},
			URL:            "",
			expectedBody:   apierror.NewBadRequestApiError(domain.ErrEmptyField),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal_error",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ResendEmailConfirmationEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ResendEmailConfirmationEmailMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			URL:            "/test?email=test",
			expectedBody:   apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ResendEmailConfirmationEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ResendEmailConfirmationEmailMockID: nil,
					},
				},
			},
			URL:            "/test?email=test",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			for _, param := range tt.params {
				c.Params = append(c.Params, param)
			}

			req, err := http.NewRequest(http.MethodGet, tt.URL, nil)
			if tt.requestBody != nil {
				req, err = http.NewRequest(http.MethodPost, tt.URL, strings.NewReader(tt.requestBody.(string)))
			}
			if err != nil {
				panic(err)
			}
			c.Request = req

			ctr := NewController(tt.fields.svc)

			ctr.ResendEmailConfirmation(c)

			response := w.Result()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("[SignUp] Expected status code = %v, got %v", tt.expectedStatus, response.StatusCode)
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			body := buf.String()

			expected, _ := json.Marshal(tt.expectedBody)
			if !reflect.DeepEqual(body, string(expected)) && tt.expectedBody != nil {
				t.Errorf("[SignUp] Expected body = %v, got %v", tt.expectedBody, body)
				return
			}
		})
	}
}

func Test_recoveryController_ConfirmEmail(t *testing.T) {
	type fields struct {
		svc RecoveryService
	}

	type args struct {
		c *gin.Context
	}

	tests := []struct {
		name           string
		fields         fields
		expectedBody   interface{}
		expectedStatus int
		URL            string
		requestBody    interface{}
		params         []gin.Param
	}{
		{
			name: "empty_email",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ConfirmEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmEmailMockID: nil,
					},
				},
			},
			URL:            "/test?token=123",
			expectedBody:   apierror.NewBadRequestApiError(domain.ErrEmptyField),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty_token",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ConfirmEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmEmailMockID: nil,
					},
				},
			},
			URL:            "/test?email=test",
			expectedBody:   apierror.NewBadRequestApiError(domain.ErrEmptyField),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal_error",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ConfirmEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmEmailMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			URL:            "/test?email=test&token=123",
			expectedBody:   apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						ConfirmEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						ConfirmEmailMockID: nil,
					},
				},
			},
			URL:            "/test?email=test&token=123",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			for _, param := range tt.params {
				c.Params = append(c.Params, param)
			}

			req, err := http.NewRequest(http.MethodGet, tt.URL, nil)
			if tt.requestBody != nil {
				req, err = http.NewRequest(http.MethodPost, tt.URL, strings.NewReader(tt.requestBody.(string)))
			}
			if err != nil {
				panic(err)
			}
			c.Request = req

			ctr := NewController(tt.fields.svc)

			ctr.ConfirmEmail(c)

			response := w.Result()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("[SignUp] Expected status code = %v, got %v", tt.expectedStatus, response.StatusCode)
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			body := buf.String()

			expected, _ := json.Marshal(tt.expectedBody)
			if !reflect.DeepEqual(body, string(expected)) && tt.expectedBody != nil {
				t.Errorf("[SignUp] Expected body = %v, got %v", tt.expectedBody, body)
				return
			}
		})
	}
}

func Test_recoveryController_SendConfirmationEmail(t *testing.T) {
	type fields struct {
		svc RecoveryService
	}

	type args struct {
		c *gin.Context
	}

	tests := []struct {
		name           string
		fields         fields
		expectedBody   interface{}
		expectedStatus int
		URL            string
		requestBody    interface{}
		params         []gin.Param
	}{
		{
			name: "empty_id",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendConfirmationEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendConfirmationEmailMockID: nil,
					},
				},
			},
			params: []gin.Param{{
				Key:   "id",
				Value: "",
			}},
			expectedBody:   apierror.NewBadRequestApiError(ErrMissingUserId),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "wrong_id",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendConfirmationEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendConfirmationEmailMockID: nil,
					},
				},
			},
			params: []gin.Param{{
				Key:   "id",
				Value: "test",
			}},
			expectedBody:   apierror.NewBadRequestApiError(ErrMissingUserId),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal_error",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendConfirmationEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendConfirmationEmailMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			params: []gin.Param{{
				Key:   "id",
				Value: "123",
			}},
			expectedBody:   apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						SendConfirmationEmailMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						SendConfirmationEmailMockID: nil,
					},
				},
			},
			params: []gin.Param{{
				Key:   "id",
				Value: "123",
			}},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			for _, param := range tt.params {
				c.Params = append(c.Params, param)
			}

			req, err := http.NewRequest(http.MethodGet, tt.URL, nil)
			if tt.requestBody != nil {
				req, err = http.NewRequest(http.MethodPost, tt.URL, strings.NewReader(tt.requestBody.(string)))
			}
			if err != nil {
				panic(err)
			}
			c.Request = req

			ctr := NewController(tt.fields.svc)

			ctr.SendConfirmationEmail(c)

			response := w.Result()

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("[SignUp] Expected status code = %v, got %v", tt.expectedStatus, response.StatusCode)
				return
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(response.Body)
			body := buf.String()

			expected, _ := json.Marshal(tt.expectedBody)
			if !reflect.DeepEqual(body, string(expected)) && tt.expectedBody != nil {
				t.Errorf("[SignUp] Expected body = %v, got %v", tt.expectedBody, body)
				return
			}
		})
	}
}
