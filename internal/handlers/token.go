package handlers

import (
	"encoding/base64"
	"go-service/internal/crypto"
	"go-service/internal/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Token token request `application/x-www-form-urlencoded` form
type TokenQuery struct {
	GrantType string `form:"grant_type,required" validate:"required"`
}

// AccessTokenResponse
// @description Access Token response body
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Token
// @Description  OAuth2 authentication only support Client Credentials grant type. required `client_id` and `client_secret` on Basic Authentication to.
// @Tags         OAuth2
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param 		 grant_type formData string true "The grant_type parameter must must be `client_credentials`" default(client_credentials)
// @Success      200  {object}  AccessTokenResponse
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /oauth/token [post]
// @Security	 BasicAuth
func (h *handler) Token(c *fiber.Ctx) error {
	tokenQuery := new(TokenQuery)
	err := c.BodyParser(tokenQuery)
	if err := h.validateStruct(tokenQuery); err != nil {
		return err
	}
	if err != nil {
		return &errorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "invalid query parameters",
		}
	}
	if tokenQuery.GrantType != "client_credentials" {
		return &errorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "grant_type only support `client_credentials`",
		}
	}

	// parse basic authentication credentials
	basic := string(c.Request().Header.Peek("Authorization"))
	if !strings.HasPrefix(basic, "Basic ") {
		return &errorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: errMessageUnauthorized,
		}
	}
	basicBase64 := basic[6:]
	basicDecode, err := base64.StdEncoding.DecodeString(basicBase64)
	if err != nil {
		return &errorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: errMessageUnauthorized,
		}
	}
	basicTokenSplit := strings.Split(string(basicDecode), ":")
	if len(basicTokenSplit) != 2 {
		return &errorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: errMessageUnauthorized,
		}
	}
	clientId := basicTokenSplit[0]
	clientSecret := basicTokenSplit[1]

	// validate client secret from database
	client := &models.Client{}
	if err := h.store.GetClientByID(client, clientId); err != nil {
		return &errorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: errMessageUnauthorized,
			Err:     err,
		}
	}
	if err := crypto.BcryptVerify(client.Secret, clientSecret, h.hashSalt); err != nil {
		return &errorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: errMessageUnauthorized,
			Err:     err,
		}
	}

	// generate JWT token
	duration := time.Hour // TODO: change duration
	token, err := h.jwt.GenerateJWT(client, duration)
	if err != nil {
		return &errorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: errMessageInternal,
			Err:     err,
		}
	}

	c.JSON(AccessTokenResponse{
		AccessToken: token,
		ExpiresIn:   int64(duration.Seconds()),
		TokenType:   "Bearer",
	})
	return nil
}
