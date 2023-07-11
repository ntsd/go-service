package crypto

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"go-service/internal/config"
	"go-service/internal/models"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTFactory interface {
	GenerateJWT(client *models.Client, duration time.Duration) (string, error)
	ValidateJWT(tokenString string) (*jwt.StandardClaims, error)
}

type jwtFactory struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewJWTFactory(cfg *config.Config) (JWTFactory, error) {
	// load JWT keys
	absPath, _ := filepath.Abs(cfg.PrivateKeyPath)
	privPEM, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privPEM)
	if err != nil {
		return nil, err
	}

	absPath, _ = filepath.Abs(cfg.PublicKeyPath)
	pubPEM, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(pubPEM)
	if err != nil {
		return nil, err
	}

	return &jwtFactory{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (f *jwtFactory) GenerateJWT(client *models.Client, duration time.Duration) (string, error) {
	// create claims
	now := time.Now()
	expireAt := now.Add(duration)
	claims := jwt.StandardClaims{
		Audience:  client.ID,
		Subject:   client.ID,
		Issuer:    "go-service",
		IssuedAt:  now.Unix(),
		ExpiresAt: expireAt.Unix(),
	}

	// create jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	if token == nil {
		return "", errors.New("failed to create jwt token")
	}

	// sign jwt token
	access, err := token.SignedString(f.privateKey)
	if err != nil {
		return "", err
	}

	return access, nil
}

func (f *jwtFactory) ValidateJWT(tokenString string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return f.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid bearer token")
	}

	return claims, nil
}
