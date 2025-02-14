package v1

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"ums/internal/dto"
	"ums/internal/logger"
	"ums/internal/mocks/v1mocks"
	"ums/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	createUserWithRolePath = "/api/v1/admin/user"
)

func TestCreateUserWithRole(t *testing.T) {
	// Initialize logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	wrappedHandler := logger.NewLogRequestHandler(handler)
	slog.SetDefault(slog.New(wrappedHandler))

	// Initialize validator
	validator, err := validation.NewValidator()
	if err != nil {
		t.Fatalf("Failed to initialize validator: %s", err)
	}

	// Valid test data
	validName := "user1"
	validPassword := "Str0ngPassword!"
	validEmail := "user1@mail.com"
	validRole := "admin"
	validUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(validName))

	validRequest := dto.CreateUserWithRoleRequest{
		Name:     validName,
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	// Invalid test data
	badNameTooShort := dto.CreateUserWithRoleRequest{
		Name:     "u1",
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	badNameTooLong := dto.CreateUserWithRoleRequest{
		Name:     "thisisaverylongusernameexceedingthelimitof64charactersaaaaaaaaaaa",
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	badNameInvalidFormat := dto.CreateUserWithRoleRequest{
		Name:     "user@1",
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	badPasswordTooShort := dto.CreateUserWithRoleRequest{
		Name:     validName,
		Email:    validEmail,
		Password: "short",
		Role:     validRole,
	}

	badPasswordNoComplexity := dto.CreateUserWithRoleRequest{
		Name:     validName,
		Email:    validEmail,
		Password: "passwordwithoutcomplexity",
		Role:     validRole,
	}

	badEmailInvalidFormat := dto.CreateUserWithRoleRequest{
		Name:     validName,
		Email:    "usermail.com",
		Password: validPassword,
		Role:     validRole,
	}

	badRoleInvalid := dto.CreateUserWithRoleRequest{
		Name:     validName,
		Email:    validEmail,
		Password: validPassword,
		Role:     "invalidRole",
	}

	missingName := dto.CreateUserWithRoleRequest{
		Name:     "",
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	missingEmail := dto.CreateUserWithRoleRequest{
		Name:     validName,
		Email:    "",
		Password: validPassword,
		Role:     validRole,
	}

	missingPassword := dto.CreateUserWithRoleRequest{
		Name:     validName,
		Email:    validEmail,
		Password: "",
		Role:     validRole,
	}

	missingRole := dto.CreateUserWithRoleRequest{
		Name:     validName,
		Email:    validEmail,
		Password: validPassword,
		Role:     "",
	}

	// Test case structure
	type args struct {
		router *gin.Engine
		url    string
		method string
		input  any
	}

	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		// Valid request
		{
			name: "Create user with role - valid request",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().CreateUserWithRole(mock.Anything, validRequest).Return(validUID, nil)
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  validRequest,
			},
			want: want{
				code:     http.StatusCreated,
				response: marshalResponse(t, dto.CreateUserResponse{ID: validUID}),
			},
		},
		{
			name: "Create user with role - internal error",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().CreateUserWithRole(mock.Anything, validRequest).Return(uuid.Nil, errors.New("service error"))
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  validRequest,
			},
			want: want{
				code:     http.StatusInternalServerError,
				response: anyString,
			},
		},

		// Invalid name cases
		{
			name: "Create user with role - name too short",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  badNameTooShort,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		{
			name: "Create user with role - name too long",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  badNameTooLong,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		{
			name: "Create user with role - invalid name format",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  badNameInvalidFormat,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		// Invalid password cases
		{
			name: "Create user with role - password too short",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  badPasswordTooShort,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		{
			name: "Create user with role - password lacking complexity",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  badPasswordNoComplexity,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		// Invalid email cases
		{
			name: "Create user with role - invalid email format",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  badEmailInvalidFormat,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		// Invalid role cases
		{
			name: "Create user with role - invalid role",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  badRoleInvalid,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		// Missing fields cases
		{
			name: "Create user with role - missing name",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  missingName,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		{
			name: "Create user with role - missing email",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  missingEmail,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		{
			name: "Create user with role - missing password",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  missingPassword,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
		{
			name: "Create user with role - missing role",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "CreateUserWithRole")
				}),
				url:    createUserWithRolePath,
				method: http.MethodPost,
				input:  missingRole,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: anyString,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			inputJSON := marshalInput(t, tt.args.input)
			req, err := http.NewRequest(tt.args.method, tt.args.url, bytes.NewReader(inputJSON))
			if err != nil {
				t.Fatalf("Failed to create HTTP request: %s", err)
			}

			tt.args.router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
			body := w.Body.String()

			if tt.want.response == anyString {
				assert.Regexp(t, regexp.MustCompile(".*"), body)
			} else {
				assert.Equal(t, tt.want.response, body)
			}
		})
	}
}
