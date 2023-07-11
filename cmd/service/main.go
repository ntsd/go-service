package main

import (
	"fmt"
	"go-service/internal/config"
	"go-service/internal/crypto"
	"go-service/internal/handlers"
	"go-service/internal/storage"
	"log"
	"runtime"
)

// @title         Go Service
// @version       1.0
// @description   This project is made for assignments for a company. To make an example of an OAuth 2.0 service that focuses on performance, maintainability, and scalability.
// @license.name  MIT
// @license.url   https://github.com/ntsd/go-service/blob/main/LICENSE
// @BasePath  	  /v1

// @securityDefinitions.basic BasicAuth
// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl /v1/oauth/token

// @externalDocs.description  Github
// @externalDocs.url          https://github.com/ntsd/go-service
func main() {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	fmt.Printf("GOMAXPROCS: %d, NumCPU: %d\n", maxProcs, numCPU)

	cfg := config.NewConfig()

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("error creating storage: %v", err)
	}

	jwtFactory, err := crypto.NewJWTFactory(cfg)
	if err != nil {
		log.Fatalf("error creating jwt factory: %v", err)
	}

	handler, err := handlers.NewHandler(cfg, storage, jwtFactory)
	if err != nil {
		log.Fatalf("error creating handler: %v", err)
	}
	handler.Serve()
}
