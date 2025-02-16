package v1

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"ums/internal/api/v1/router"
	"ums/internal/dto"
	"ums/internal/logger"
	"ums/internal/mocks/v1mocks"
	"ums/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadUser(t *testing.T) {
	// Initialize logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	wrappedHandler := logger.NewLogRequestHandler(handler)
	slog.SetDefault(slog.New(wrappedHandler))

	// Valid test data
	validName := "user1"
	validEmail := "user1@mail.com"
	validRole := "admin"
	validUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(validName))
	now := time.Now()

	validResponse := dto.ReadUserResponse{
		ID:        validUID,
		Name:      validName,
		Email:     validEmail,
		Role:      validRole,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test case structure
	type args struct {
		router *gin.Engine
		url    string
		method string
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
			name: "Read user - valid request",
			args: args{
				router: prepareRouter(t, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().ReadUser(mock.Anything, validUID).Return(validResponse, nil)
				}),
				url:    "/api/v1/user/" + validUID.String(),
				method: http.MethodGet,
			},
			want: want{
				code:     http.StatusOK,
				response: serializeResponse(t, validResponse),
			},
		},
		// Service error case
		{
			name: "Read user - service error",
			args: args{
				router: prepareRouter(t, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().ReadUser(mock.Anything, validUID).Return(dto.ReadUserResponse{}, errors.New("service error"))
				}),
				url:    "/api/v1/user/" + validUID.String(),
				method: http.MethodGet,
			},
			want: want{
				code:     http.StatusInternalServerError,
				response: "",
			},
		},
		// Invalid UUID case
		{
			name: "Read user - invalid UUID",
			args: args{
				router: prepareRouter(t, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "ReadUser")
				}),
				url:    "/api/v1/user/invalid-uuid",
				method: http.MethodGet,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest(tt.args.method, tt.args.url, nil)
			if err != nil {
				t.Fatalf("Failed to create HTTP request: %s", err)
			}

			tt.args.router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.code, w.Code)
			body := w.Body.String()

			if tt.want.response != "" {
				assert.JSONEq(t, tt.want.response, body)
			}
		})
	}
}

func prepareRouter(t *testing.T, mockSetup func(*v1mocks.MockUMS)) *gin.Engine {
	mockums := v1mocks.NewMockUMS(t)
	mockSetup(mockums)

	// Initialize validator
	validator, err := validation.NewValidator()
	if err != nil {
		t.Fatalf("Failed to initialize validator: %s", err)
	}

	h, err := New(mockums, validator) // Assuming validator is not needed for this handler
	if err != nil {
		t.Fatalf("Failed to initialize handler: %s", err)
	}

	return router.SetupRouter(h, gin.TestMode)
}

func serializeResponse(t *testing.T, response any) string {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to serialize response: %s", err)
	}
	return string(jsonResponse)
}
