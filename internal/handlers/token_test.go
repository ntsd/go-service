package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"go-service/internal/crypto"
	mock_crypto "go-service/internal/crypto/mock_jwt"
	"go-service/internal/models"
	"go-service/internal/storage"
	"go-service/internal/storage/mock_storage"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"gorm.io/gorm"
)

func Test_handler_Token(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	tests := []testData{
		{
			name: "success",
			request: func() *http.Request {
				data := url.Values{}
				data.Set("grant_type", "client_credentials")

				req := httptest.NewRequest("POST", "/v1/oauth/token", strings.NewReader(data.Encode()))
				b64 := []byte(fmt.Sprintf("%s:%s", "string", "string"))
				req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(b64))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			mockJWTFactory: func(ctrl *gomock.Controller) crypto.JWTFactory {
				mockJWT := mock_crypto.NewMockJWTFactory(ctrl)
				mockJWT.EXPECT().GenerateJWT(gomock.Any(), gomock.Any()).Return("token", nil)
				return mockJWT
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetClientByID(gomock.Any(), "string").Return(nil).SetArg(0, models.Client{
					ID:     "string",
					Secret: []byte("$2a$10$xsytJkE2thJaWw3Sgh5Rbup1YahdKVKmX.0mYj9SF0w62XKIMYBzS"),
				})
				return mockStorage
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "not found",
			request: func() *http.Request {
				data := url.Values{}
				data.Set("grant_type", "client_credentials")

				req := httptest.NewRequest("POST", "/v1/oauth/token", strings.NewReader(data.Encode()))
				b64 := []byte(fmt.Sprintf("%s:%s", "string", "string"))
				req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(b64))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetClientByID(gomock.Any(), "string").Return(gorm.ErrRecordNotFound)
				return mockStorage
			},
			wantStatus: fiber.StatusUnauthorized,
		},
		{
			name: "wrong password",
			request: func() *http.Request {
				data := url.Values{}
				data.Set("grant_type", "client_credentials")

				req := httptest.NewRequest("POST", "/v1/oauth/token", strings.NewReader(data.Encode()))
				b64 := []byte(fmt.Sprintf("%s:%s", "string", "xxx"))
				req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(b64))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetClientByID(gomock.Any(), "string").Return(nil).SetArg(0, models.Client{
					ID:     "string",
					Secret: []byte("$2a$10$xsytJkE2thJaWw3Sgh5Rbup1YahdKVKmX.0mYj9SF0w62XKIMYBzS"),
				})
				return mockStorage
			},
			wantStatus: fiber.StatusUnauthorized,
		},
		{
			name: "error generate jwt",
			request: func() *http.Request {
				data := url.Values{}
				data.Set("grant_type", "client_credentials")

				req := httptest.NewRequest("POST", "/v1/oauth/token", strings.NewReader(data.Encode()))
				b64 := []byte(fmt.Sprintf("%s:%s", "string", "string"))
				req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(b64))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			mockJWTFactory: func(ctrl *gomock.Controller) crypto.JWTFactory {
				mockJWT := mock_crypto.NewMockJWTFactory(ctrl)
				mockJWT.EXPECT().GenerateJWT(gomock.Any(), gomock.Any()).Return("", errors.New("error"))
				return mockJWT
			},
			mockStorage: func(ctrl *gomock.Controller) storage.Storage {
				mockStorage := mock_storage.NewMockStorage(ctrl)
				mockStorage.EXPECT().GetClientByID(gomock.Any(), "string").Return(nil).SetArg(0, models.Client{
					ID:     "string",
					Secret: []byte("$2a$10$xsytJkE2thJaWw3Sgh5Rbup1YahdKVKmX.0mYj9SF0w62XKIMYBzS"),
				})
				return mockStorage
			},
			wantStatus: fiber.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h, app := mockHandlerAndTest(test, mockCtrl)

			app.Post("/v1/oauth/token", h.Token)

			resp, _ := app.Test(test.request())
			validateResponse(t, test, resp.StatusCode)
		})
	}
}
