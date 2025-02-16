package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"ums/internal/dto"
	"ums/internal/logger"
	"ums/internal/mocks/v1mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadUsers(t *testing.T) {
	// Initialize logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	wrappedHandler := logger.NewLogRequestHandler(handler)
	slog.SetDefault(slog.New(wrappedHandler))

	// Valid test data
	validName1 := "user1"
	validEmail1 := "user1@mail.com"
	validRole1 := "admin"
	validUID1 := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(validName1))
	now := time.Now()

	validName2 := "user2"
	validEmail2 := "user2@mail.com"
	validRole2 := "user"
	validUID2 := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(validName2))

	validResponse := dto.ReadUsersResponse{
		Users: []dto.ReadUserResponse{
			{
				ID:        validUID1,
				Name:      validName1,
				Email:     validEmail1,
				Role:      validRole1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        validUID2,
				Name:      validName2,
				Email:     validEmail2,
				Role:      validRole2,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
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
			name: "Read users - valid request",
			args: args{
				router: prepareRouter(t, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().ReadUsers(mock.Anything).Return(validResponse, nil)
				}),
				url:    "/api/v1/admin/users",
				method: http.MethodGet,
			},
			want: want{
				code:     http.StatusOK,
				response: serializeResponse(t, validResponse),
			},
		},
		// Service error case
		{
			name: "Read users - service error",
			args: args{
				router: prepareRouter(t, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().ReadUsers(mock.Anything).Return(dto.ReadUsersResponse{}, errors.New("service error"))
				}),
				url:    "/api/v1/admin/users",
				method: http.MethodGet,
			},
			want: want{
				code:     http.StatusInternalServerError,
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
