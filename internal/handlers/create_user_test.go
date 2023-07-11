package handlers

import (
	"errors"
	"go-service/internal/storage"
	"go-service/internal/storage/mock_storage"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Test_handler_CreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	body := `{ "email": "example@example.com", "name": "John Doe" }`

	tests := []testData{
		{
			name: "success",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/v1/users", nil)
				req.Body = io.NopCloser(strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Length", strconv.Itoa(len(body)))
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().CreateUser(gomock.Any()).Return(nil)
				return mockStorage
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "validate error no email",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/v1/users", nil)
				body := `{ "name": "John Doe" }`
				req.Body = io.NopCloser(strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Length", strconv.Itoa(len(body)))
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				return mockStorage
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "validate error email",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/v1/users", nil)
				body := `{ "email": "x", "name": "John Doe" }`
				req.Body = io.NopCloser(strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Length", strconv.Itoa(len(body)))
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				return mockStorage
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "validate error name",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/v1/users", nil)
				body := `{ "email": "example@example.com", "name": "x" }`
				req.Body = io.NopCloser(strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Length", strconv.Itoa(len(body)))
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				return mockStorage
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "database error",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/v1/users", nil)
				req.Body = io.NopCloser(strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Length", strconv.Itoa(len(body)))
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().CreateUser(gomock.Any()).Return(errors.New("error"))
				return mockStorage
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name: "pgx error",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/v1/users", nil)
				req.Body = io.NopCloser(strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Length", strconv.Itoa(len(body)))
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().CreateUser(gomock.Any()).Return(&pgconn.PgError{Code: "1"})
				return mockStorage
			},
			wantStatus: fiber.StatusUnprocessableEntity,
		},
		{
			name: "duplicate error",
			request: func() *http.Request {
				req := httptest.NewRequest("POST", "/v1/users", nil)
				req.Body = io.NopCloser(strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Length", strconv.Itoa(len(body)))
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().CreateUser(gomock.Any()).Return(&pgconn.PgError{Code: pgerrcode.UniqueViolation})
				return mockStorage
			},
			wantStatus: fiber.StatusUnprocessableEntity,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h, app := mockHandlerAndTest(test, mockCtrl)

			app.Post("/v1/users", h.CreateUser)

			resp, err := app.Test(test.request())
			if err != nil {
				t.Fatal(err)
			}
			validateResponse(t, test, resp.StatusCode)
		})
	}
}
