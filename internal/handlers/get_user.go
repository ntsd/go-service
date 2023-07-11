package handlers

import (
	"go-service/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetUser
// @Description  Get user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /users/{id} [get]
// @Security	 OAuth2Application
func (h *handler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user := &models.User{}
	if err := h.store.GetUserByID(user, id); err != nil {
		if err := dbErrorParse(err); err != nil {
			return err
		}

		return &errorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     err,
		}
	}

	c.JSON(user)
	return nil
}
