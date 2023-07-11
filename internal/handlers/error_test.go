package handlers

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func TestDbErrorParse(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name    string
		err     error
		wantErr *errorResponse
	}{
		{
			name: "gorm.ErrRecordNotFound",
			err:  gorm.ErrRecordNotFound,
			wantErr: &errorResponse{
				Code:    fiber.StatusNotFound,
				Message: errMessageNotFound,
			},
		},
		{
			name: "pgconn.PgError with UniqueViolation",
			err: &pgconn.PgError{
				Code:           pgerrcode.UniqueViolation,
				ConstraintName: "test_constraint",
			},
			wantErr: &errorResponse{
				Code:    fiber.StatusUnprocessableEntity,
				Message: "test_constraint is duplicate",
			},
		},
		{
			name: "pgconn.PgError with other error",
			err: &pgconn.PgError{
				Code:   "other_error",
				Detail: "other detail",
			},
			wantErr: &errorResponse{
				Code:    fiber.StatusUnprocessableEntity,
				Message: "other detail",
			},
		},
		{
			name:    "other error",
			err:     errors.New("other error"),
			wantErr: nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := dbErrorParse(tc.err)
			if diff := cmp.Diff(
				err,
				tc.wantErr,
				cmpopts.IgnoreFields(errorResponse{}, "Err"),
			); diff != "" {
				t.Fatalf("want error mismatch(-want +got):\n%s", diff)
			}
		})
	}
}
