package handlers

import (
	"errors"
	"go-service/internal/storage"
	"go-service/internal/storage/mock_storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
)

func Test_handler_GetUsers(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	tests := []testData{
		{
			name: "success",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListUsers(gomock.Any(), 0, 100, "").Return(nil)
				return mockStorage
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "success change offset",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users?offset=10", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListUsers(gomock.Any(), 10, 100, "").Return(nil)
				return mockStorage
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "success change limit",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users?limit=10", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListUsers(gomock.Any(), 0, 10, "").Return(nil)
				return mockStorage
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "success change name",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users?name=test", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListUsers(gomock.Any(), 0, 100, "test").Return(nil)
				return mockStorage
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "error offset not number",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users?offset=test", nil)
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "error limit not number",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users?limit=test", nil)
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "database error",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListUsers(gomock.Any(), 0, 100, "").Return(errors.New("error"))
				return mockStorage
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name: "pgx error",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().ListUsers(gomock.Any(), 0, 100, "").Return(&pgconn.PgError{Code: "1"})
				return mockStorage
			},
			wantStatus: fiber.StatusUnprocessableEntity,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h, app := mockHandlerAndTest(test, mockCtrl)

			app.Get("/v1/users", h.GetUsers)

			resp, _ := app.Test(test.request())
			validateResponse(t, test, resp.StatusCode)
		})
	}
}
