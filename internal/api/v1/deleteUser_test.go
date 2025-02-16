package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"ums/internal/logger"
	"ums/internal/mocks/v1mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUser(t *testing.T) {
	// Initialize logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	wrappedHandler := logger.NewLogRequestHandler(handler)
	slog.SetDefault(slog.New(wrappedHandler))

	// Valid test data
	validName := "user1"
	validUID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(validName))

	// Test case structure
	type args struct {
		router *gin.Engine
		url    string
		method string
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
			name: "Delete user - valid request",
			args: args{
				router: prepareRouter(t, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().DeleteUser(mock.Anything, validUID).Return(nil)
				}),
				url:    "/api/v1/user/" + validUID.String(),
				method: http.MethodDelete,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		// Service error case
		{
			name: "Delete user - service error",
			args: args{
				router: prepareRouter(t, func(mockums *v1mocks.MockUMS) {
					mockums.EXPECT().DeleteUser(mock.Anything, validUID).Return(errors.New("service error"))
				}),
				url:    "/api/v1/user/" + validUID.String(),
				method: http.MethodDelete,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		// Invalid UUID case
		{
			name: "Delete user - invalid UUID",
			args: args{
				router: prepareRouter(t, func(mockums *v1mocks.MockUMS) {
					mockums.AssertNotCalled(t, "DeleteUser")
				}),
				url:    "/api/v1/user/invalid-uuid",
				method: http.MethodDelete,
			},
			want: want{
				code: http.StatusBadRequest,
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
		})
	}
}
