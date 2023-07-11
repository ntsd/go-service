package handlers

import (
	"go-service/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateUserBody create user body
type CreateUserBody struct {
	Name  string `json:"name" validate:"required,min=3,max=32" default:"John Doe"`
	Email string `json:"email" validate:"required,email,min=6,max=32" default:"example@example.com"`
}

// CreateUser
// @Description  create a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param 		 user body CreateUserBody true "JSON body of the user"
// @Success      200  {object}  models.User
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      422  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /users [post]
// @Security	 OAuth2Application
func (h *handler) CreateUser(c *fiber.Ctx) error {
	userBody := new(CreateUserBody)
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
		ID:    uuid.New().String(),
		Email: userBody.Email,
		Name:  userBody.Name,
	}
	if err := h.store.CreateUser(newUser); err != nil {
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
