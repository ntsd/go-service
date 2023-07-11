package handlers

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	errMessageBadRequest   string = "bad request"
	errMessageNotFound     string = "not found"
	errMessageDuplicate    string = "duplicate"
	errMessageInternal     string = "something went wrong"
	errMessageUnauthorized string = "unauthorized"
)

// errorResponse use for response the error
// @Description Common error response.
type errorResponse struct {
	// Code HTTP status code
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error is error implement of the standard error interface
func (e *errorResponse) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// errorHandler Fiber error handler
func errorHandler(c *fiber.Ctx, err error) error {
	log.Println(err)

	// set Content-Type: application/json; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	// retrieve the custom status code if it's a *errorResponse
	var e *errorResponse
	if errors.As(err, &e) {
		// return status code with error message
		return c.Status(e.Code).JSON(e)
	}

	// return status code 500 with error message
	return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{
		Message: "Internal Server Error",
	})
}

// dbErrorParse parse pgx/gorm error to response error
func dbErrorParse(err error) *errorResponse {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return &errorResponse{
			Code:    fiber.StatusNotFound,
			Message: errMessageNotFound,
			Err:     err,
		}
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return &errorResponse{
				Code:    fiber.StatusUnprocessableEntity,
				Message: fmt.Sprintf("%s is duplicate", pgErr.ConstraintName),
				Err:     err,
			}
		default:
			return &errorResponse{
				Code:    fiber.StatusUnprocessableEntity,
				Message: pgErr.Detail,
				Err:     err,
			}
		}
	}

	return nil
}
