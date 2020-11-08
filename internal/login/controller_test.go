package login

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
	LoginUserMockID = iota
	UserCanLoginMockID
)

type MockService struct {
	Responses map[int]interface{}
	Errors    map[int]apierror.ApiError
}

func (m *MockService) LoginUser(user *domain.UserLoginDTO, ctx *rest.ContextInformation) (string, apierror.ApiError) {
	return m.Responses[LoginUserMockID].(string), m.Errors[LoginUserMockID]
}

func (m *MockService) UserCanLogin(user *domain.UserLoginDTO) apierror.ApiError {
	return m.Errors[UserCanLoginMockID]
}

func Test_loginController_Login(t *testing.T) {
	type fields struct {
		svc Service
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
			name: "bad_request",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						LoginUserMockID: true,
					},
					Errors: map[int]apierror.ApiError{
						LoginUserMockID: nil,
					},
				},
			},
			requestBody:    `"}`,
			expectedBody:   apierror.New(http.StatusBadRequest, domain.ErrInvalidBody, apierror.NewErrorCause(domain.ErrInvalidBody, domain.ErrInvalidBodyCode)),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal_error",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						LoginUserMockID: "test",
					},
					Errors: map[int]apierror.ApiError{
						LoginUserMockID: apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
					},
				},
			},
			requestBody:    `{"user": "test"}`,
			expectedBody:   apierror.NewInternalServerApiError("Internal error", errors.New("error"), "test"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ok",
			fields: fields{
				svc: &MockService{
					Responses: map[int]interface{}{
						LoginUserMockID: "test",
					},
					Errors: map[int]apierror.ApiError{
						LoginUserMockID: nil,
					},
				},
			},
			requestBody:    `{"user": "test"}`,
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

			ctr.Login(c)

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
