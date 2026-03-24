package service

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/storage"
)

type TokenService struct {
	parser *jwt.Parser

	storage *storage.TokenStorage

	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func NewTokenService(storage *storage.TokenStorage) (*TokenService, error) {
	privateKeyFile, err := os.ReadFile("/data/keys/ed25519-priv.pem")
	if err != nil {
		return nil, fmt.Errorf("unable to read private key file: %w", err)
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	ed25519PrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not Ed25519")
	}

	ed25519PublicKey, ok := ed25519PrivateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not Ed25519")
	}

	signingMethod := jwt.SigningMethodEdDSA

	return &TokenService{
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{signingMethod.Alg()}),
			jwt.WithIssuer("oxyl"),
			jwt.WithIssuedAt(),
			jwt.WithExpirationRequired(),
			jwt.WithLeeway(15*time.Second),
		),
		storage:    storage,
		publicKey:  ed25519PublicKey,
		privateKey: ed25519PrivateKey,
	}, nil
}

func (t *TokenService) CreateToken(identifier string, holder *string, tokenType models.JWTTokenType) (*models.TokenPair, error) {
	accessToken, err := models.NewToken(identifier, holder, tokenType)
	if err != nil {
		return nil, fmt.Errorf("unable to create access token: %w", err)
	}

	refreshToken, err := models.NewRefreshToken(accessToken.Identifier, accessToken.ID, tokenType)
	if err != nil {
		return nil, fmt.Errorf("unable to create refresh token: %w", err)
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, accessToken.RegisteredClaims).SignedString(t.privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to sign access token: %w", err)
	}

	refreshTokenString, err := jwt.NewWithClaims(jwt.SigningMethodEdDSA, refreshToken.RegisteredClaims).SignedString(t.privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to sign refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken: struct {
			Token     string    `json:"token"`
			ExpiresAt time.Time `json:"expires_at"`
		}{
			Token:     tokenString,
			ExpiresAt: accessToken.ExpiresAt.Time,
		},
		RefreshToken: struct {
			Token     string    `json:"token"`
			ExpiresAt time.Time `json:"expires_at"`
		}{
			Token:     refreshTokenString,
			ExpiresAt: refreshToken.ExpiresAt.Time,
		},
	}, nil
}

func (t *TokenService) ParseToken(token string) (*models.Token, error) {
	tokenParsed, err := t.parser.ParseWithClaims(token, &models.Token{}, func(token *jwt.Token) (any, error) {
		return t.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	claims, ok := tokenParsed.Claims.(*models.Token)
	if !ok {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	return claims, nil

}

func (t *TokenService) ParseRefreshToken(token string) (*models.RefreshToken, error) {
	tokenParsed, err := t.parser.ParseWithClaims(token, &models.RefreshToken{}, func(token *jwt.Token) (any, error) {
		return t.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse refresh token: %w", err)
	}

	claims, ok := tokenParsed.Claims.(*models.RefreshToken)
	if !ok {
		return nil, fmt.Errorf("unable to parse refresh token: %w", err)
	}

	return claims, nil
}

func (t *TokenService) RefreshToken(refreshToken string) (string, error) {
	return "", nil
}

func (t *TokenService) RevokeToken(ctx context.Context, refreshToken string) error {
	return nil
}
