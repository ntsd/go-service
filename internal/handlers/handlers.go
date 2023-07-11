package handlers

import (
	"fmt"
	"go-service/internal/config"
	"go-service/internal/crypto"
	"go-service/internal/storage"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	_ "go-service/docs"

	"github.com/go-playground/validator"
)

// Handler is the http server handlers interface
type Handler interface {
	CreateClient(c *fiber.Ctx) error
	Token(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error

	Serve()
}

// handler is the http server handlers struct
type handler struct {
	address  string
	store    storage.Storage
	jwt      crypto.JWTFactory
	hashSalt string
	preFork  bool

	validator *validator.Validate
}

// NewHandler create a new handler
func NewHandler(cfg *config.Config, store storage.Storage, jwt crypto.JWTFactory) (Handler, error) {
	return &handler{
		hashSalt:  cfg.HashSalt,
		address:   fmt.Sprintf(":%s", cfg.AppPort),
		store:     store,
		jwt:       jwt,
		validator: validator.New(),
	}, nil
}

// Serve start the http server
func (h *handler) Serve() {
	app := fiber.New(fiber.Config{
		Prefork:      h.preFork,
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())

	// Swagger handlers
	app.Get("/swagger/*", swagger.HandlerDefault)

	// OAuth2 handlers
	app.Post("/v1/oauth/clients", h.CreateClient)
	app.Post("/v1/oauth/token", h.Token)

	// JWT Middleware
	app.Use(h.JWTMiddleware)

	// Users handlers
	app.Get("/v1/users", h.GetUsers)
	app.Post("/v1/users", h.CreateUser)
	app.Get("/v1/users/:id", h.GetUser)
	app.Put("/v1/users/:id", h.UpdateUser)

	log.Fatal(app.Listen(h.address))
}
