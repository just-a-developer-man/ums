package v1

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestUpdateUserData(t *testing.T) {
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

	validRequest := dto.UpdateUserDataRequest{
		Name:     validName,
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	// Invalid test data
	badNameInvalidFormat := dto.UpdateUserDataRequest{
		Name:     "user@1",
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	badPasswordTooShort := dto.UpdateUserDataRequest{
		Name:     validName,
		Email:    validEmail,
		Password: "short",
		Role:     validRole,
	}

	badPasswordNoComplexity := dto.UpdateUserDataRequest{
		Name:     validName,
		Email:    validEmail,
		Password: "passwordwithoutcomplexity",
		Role:     validRole,
	}

	badEmailInvalidFormat := dto.UpdateUserDataRequest{
		Name:     validName,
		Email:    "usermail.com",
		Password: validPassword,
		Role:     validRole,
	}

	badRoleInvalid := dto.UpdateUserDataRequest{
		Name:     validName,
		Email:    validEmail,
		Password: validPassword,
		Role:     "invalidRole",
	}

	missingName := dto.UpdateUserDataRequest{
		Name:     "",
		Email:    validEmail,
		Password: validPassword,
		Role:     validRole,
	}

	missingEmail := dto.UpdateUserDataRequest{
		Name:     validName,
		Email:    "",
		Password: validPassword,
		Role:     validRole,
	}

	missingPassword := dto.UpdateUserDataRequest{
		Name:     validName,
		Email:    validEmail,
		Password: "",
		Role:     validRole,
	}

	missingRole := dto.UpdateUserDataRequest{
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
		code int
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		// Valid request
		{
			name: "Update user data - valid request",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().UpdateUserData(mock.Anything, validUID, &validRequest).Return(nil)
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  validRequest,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		// Service error case
		{
			name: "Update user data - service error",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().UpdateUserData(mock.Anything, validUID, &validRequest).Return(errors.New("service error"))
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  validRequest,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		// Invalid name cases
		{
			name: "Update user data - invalid name format",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  badNameInvalidFormat,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		// Invalid password cases
		{
			name: "Update user data - password too short",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  badPasswordTooShort,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Update user data - password lacking complexity",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  badPasswordNoComplexity,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		// Invalid email cases
		{
			name: "Update user data - invalid email format",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  badEmailInvalidFormat,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		// Invalid role cases
		{
			name: "Update user data - invalid role",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  badRoleInvalid,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		// Missing fields cases
		{
			name: "Update user data - missing name",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  missingName,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Update user data - missing email",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  missingEmail,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Update user data - missing password",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  missingPassword,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Update user data - missing role",
			args: args{
				router: setupRouter(t, validator, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "UpdateUserData")
				}),
				url:    "/api/v1/admin/" + validUID.String(),
				method: http.MethodPut,
				input:  missingRole,
			},
			want: want{
				code: http.StatusBadRequest,
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
		})
	}
}
