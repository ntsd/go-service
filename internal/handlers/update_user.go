package handlers

import (
	"go-service/internal/models"

	"github.com/gofiber/fiber/v2"
)

// UpdateUserBody update user body
type UpdateUserBody struct {
	Name  string `json:"name" validate:"required,min=3,max=32" default:"John Doe"`
	Email string `json:"email" validate:"required,email,min=6,max=32" default:"example@example.com"`
}

// UpdateUser
// @Description  update a user, it will not insert if not existing.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param 		 user body UpdateUserBody true "JSON body of the user"
// @Success      200  {object}  models.User
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      422  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /users/{id} [put]
// @Security	 OAuth2Application
func (h *handler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	userBody := new(UpdateUserBody)
	err := c.BodyParser(userBody)
	if err := h.validateStruct(userBody); err != nil {
		return err
	}
	if err != nil {
		return &errorResponse{
			Code:    fiber.StatusBadRequest,
			Message: errMessageBadRequest,
			Err:     err,
		}
	}

	newUser := &models.User{
		Email: userBody.Email,
		Name:  userBody.Name,
	}
	if err := h.store.UpdateUserByID(newUser, id); err != nil {
		if err := dbErrorParse(err); err != nil {
			return err
		}

		return &errorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     err,
		}
	}

	c.JSON(newUser)
	return nil
}
