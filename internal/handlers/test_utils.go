package handlers

import (
	"go-service/internal/crypto"
	"go-service/internal/models"
	"go-service/internal/storage"
	"net/http"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// testData struct for handler test
type testData struct {
	name           string
	request        func() *http.Request
	mockStorage    func(ctrl *gomock.Controller) storage.Storage
	mockJWTFactory func(ctrl *gomock.Controller) crypto.JWTFactory
	wantStatus     int
}

// mockHandlerAndTest create mock handler and gin context
func mockHandlerAndTest(test testData, mockCtrl *gomock.Controller) (*handler, *fiber.App) {
	// initial handler function and mock
	h := &handler{
		validator: validator.New(),
		hashSalt:  "change_this_salt",
	}
	if test.mockJWTFactory != nil {
		h.jwt = test.mockJWTFactory(mockCtrl)
	}
	if test.mockStorage != nil {
		h.store = test.mockStorage(mockCtrl)
	}

	app := fiber.New(fiber.Config{
		Prefork:      false,
		ErrorHandler: errorHandler,
	})

	return h, app
}

// validateResponse check diff data response and error response, fatal when unmatched
func validateResponse(
	t *testing.T,
	test testData,
	status int,
) {
	// validate data response
	if diff := cmp.Diff(
		test.wantStatus,
		status,
		cmpopts.IgnoreFields(models.User{}, "ID"),
	); diff != "" {
		t.Fatalf("want data mismatch(-want +got):\n%s", diff)
	}
}
