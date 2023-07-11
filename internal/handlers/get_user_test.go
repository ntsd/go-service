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
	"gorm.io/gorm"
)

func Test_handler_GetUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	tests := []testData{
		{
			name: "success",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users/foo", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetUserByID(gomock.Any(), "foo").Return(nil)
				return mockStorage
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "not found",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users/foo", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetUserByID(gomock.Any(), "foo").Return(gorm.ErrRecordNotFound)
				return mockStorage
			},
			wantStatus: fiber.StatusNotFound,
		},
		{
			name: "database error",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users/foo", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetUserByID(gomock.Any(), "foo").Return(errors.New("error"))
				return mockStorage
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name: "pgx error",
			request: func() *http.Request {
				return httptest.NewRequest("GET", "/v1/users/foo", nil)
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetUserByID(gomock.Any(), "foo").Return(&pgconn.PgError{Code: "1"})
				return mockStorage
			},
			wantStatus: fiber.StatusUnprocessableEntity,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h, app := mockHandlerAndTest(test, mockCtrl)

			app.Get("/v1/users/:id", h.GetUser)

			resp, _ := app.Test(test.request())
			validateResponse(t, test, resp.StatusCode)
		})
	}
}
