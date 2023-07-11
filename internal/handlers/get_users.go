package handlers

import (
	"go-service/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetUsersQuery get users request query parameters
type GetUsersQuery struct {
	Offset int    `query:"offset" validate:"number"`
	Limit  int    `query:"limit" validate:"number"`
	Name   string `query:"name"`
}

// GetUsers
// @Description  List users.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param 		 offset query int false "pagination offset, default is `0`"
// @Param 		 limit query int false "pagination limit, default is `100`. If more than `100` will be set as `100`."
// @Param 		 name query string false "filter name by partial text search"
// @Success      200  {array}   models.User
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /users [get]
// @Security	 OAuth2Application
func (h *handler) GetUsers(c *fiber.Ctx) error {
	userQuery := new(GetUsersQuery)
	err := c.QueryParser(userQuery)
	if err := h.validateStruct(userQuery); err != nil {
		return err
	}
	if err != nil {
		return &errorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid query parameters",
		}
	}
	if userQuery.Limit == 0 || userQuery.Limit > 100 {
		userQuery.Limit = 100
	}

	users := &[]models.User{}
	if err := h.store.ListUsers(users, userQuery.Offset, userQuery.Limit, userQuery.Name); err != nil {
		if err := dbErrorParse(err); err != nil {
			return err
		}

		return &errorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     err,
		}
	}

	c.JSON(users)
	return nil
}
