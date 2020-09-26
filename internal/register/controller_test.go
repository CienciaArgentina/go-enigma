package register

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
	UserCanSignUpMockName = "UserCanSignUp"
	CreateUserMockName    = "CreateUser"
)

type MockService struct {
	Responses map[string]interface{}
	Errors    map[string]apierror.ApiError
}

func (m *MockService) UserCanSignUp(u *domain.UserSignupDTO) (bool, apierror.ApiError) {
	return m.Responses[UserCanSignUpMockName].(bool), m.Errors[UserCanSignUpMockName]
}

func (m *MockService) CreateUser(u *domain.UserSignupDTO, ctx *rest.ContextInformation) (int64, apierror.ApiError) {
	return m.Responses[CreateUserMockName].(int64), m.Errors[CreateUserMockName]
}

func Test_registerController_SignUp(t *testing.T) {
	type fields struct {
		svc RegisterService
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
			name:           "invalid_request",
			expectedStatus: http.StatusBadRequest,
			requestBody:    "{2",
			expectedBody:   apierror.New(http.StatusBadRequest, domain.ErrInvalidBody, apierror.NewErrorCause(domain.ErrInvalidBody, domain.ErrInvalidBodyCode)),
		},
		{
			name: "service_error",
			fields: fields{
				svc: &MockService{
					Responses: map[string]interface{}{
						CreateUserMockName: int64(0),
					},
					Errors: map[string]apierror.ApiError{
						CreateUserMockName: apierror.NewInternalServerApiError("Error! :(", errors.New(":("), "test"),
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
			requestBody:    "{}",
			expectedBody:   apierror.NewInternalServerApiError("Error! :(", errors.New(":("), "test"),
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[string]interface{}{
						CreateUserMockName: int64(123),
					},
				},
			},
			expectedStatus: http.StatusOK,
			requestBody:    "{}",
			expectedBody:   gin.H{"user_id": 123},
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

			ctr.SignUp(c)

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
