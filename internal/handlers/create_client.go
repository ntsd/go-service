package handlers

import (
	"go-service/internal/crypto"
	"go-service/internal/models"

	"github.com/gofiber/fiber/v2"
)

// CreateClientBody create client body
type CreateClientBody struct {
	ClientID     string `json:"client_id" validate:"required,max=64"`
	ClientSecret string `json:"client_secret" validate:"required,max=64"`
}

// CreateClient
// @Description  create a new client
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Param 		 client body CreateClientBody true "JSON body of the client"
// @Success      200  {object}  CreateClientBody
// @Failure      400  {object}  errorResponse
// @Failure      422  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /oauth/clients [post]
func (h *handler) CreateClient(c *fiber.Ctx) error {
	clientBody := new(CreateClientBody)
	err := c.BodyParser(clientBody)
	if err := h.validateStruct(clientBody); err != nil {
		return err
	}
	if err != nil {
		return &errorResponse{
			Code:    fiber.StatusBadRequest,
			Message: errMessageBadRequest,
			Err:     err,
		}
	}

	hash, err := crypto.BcryptHash(clientBody.ClientSecret, h.hashSalt)
	if err != nil {
		return &errorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     err,
		}
	}

	client := &models.Client{
		ID:     clientBody.ClientID,
		Secret: hash,
	}
	if err := h.store.CreateClient(client); err != nil {
		if err := dbErrorParse(err); err != nil {
			return err
		}

		return &errorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     err,
		}
	}

	c.JSON(map[string]string{
		"id": clientBody.ClientID,
	})
	return nil
}
