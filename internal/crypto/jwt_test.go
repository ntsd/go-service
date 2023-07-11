package crypto

import (
	"go-service/internal/config"
	"go-service/internal/models"
	"testing"
	"time"
)

func TestJWTFactory(t *testing.T) {
	cfg := &config.Config{
		PrivateKeyPath: "../../deployments/ec_private.pem",
		PublicKeyPath:  "../../deployments/ec_public.pem",
	}
	factory, err := NewJWTFactory(cfg)
	if err != nil {
		t.Fatalf("Failed to create JWTFactory: %v", err)
	}

	client := &models.Client{ID: "test_client"}
	duration := time.Hour
	token, err := factory.GenerateJWT(client, duration)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	claims, err := factory.ValidateJWT(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if claims.Audience != client.ID {
		t.Errorf("Expected Audience %v, got %v", client.ID, claims.Audience)
	}
	if claims.Subject != client.ID {
		t.Errorf("Expected Subject %v, got %v", client.ID, claims.Subject)
	}
	if claims.Issuer != "go-service" {
		t.Errorf("Expected Issuer %v, got %v", "go-service", claims.Issuer)
	}
}
