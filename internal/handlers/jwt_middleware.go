package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware fiber middleware
func (h *handler) JWTMiddleware(c *fiber.Ctx) error {
	// parse Bearer token
	bearer := string(c.Request().Header.Peek("Authorization"))
	if !strings.HasPrefix(bearer, "Bearer ") {
		return &errorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: errMessageUnauthorized,
		}
	}
	token := bearer[7:]

	if _, err := h.jwt.ValidateJWT(token); err != nil {
		return &errorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: errMessageUnauthorized,
			Err:     err,
		}
	}

	return c.Next()
}
