package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type Token struct {
	Identifier string `json:"identifier"`

	// ? this can be nil if the token is not bound to a specific company, and it is not from a user
	Holder *string      `json:"holder,omitempty"`
	Type   JWTTokenType `json:"type"`

	jwt.RegisteredClaims
}

func NewToken(identifier string, holder *string, tokenType JWTTokenType) (*Token, error) {
	if identifier == "" {
		return nil, errors.New("identifier is empty")
	}

	if tokenType != TokenTypeAgent && tokenType != TokenTypeUser {
		return nil, errors.New("token type is not valid")
	}

	if tokenType == TokenTypeAgent && (holder == nil || *holder == "") {
		return nil, errors.New("holder is nil")
	}

	return new(Token{
		Holder:     holder,
		Identifier: identifier,
		Type:       tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        ulid.Make().String(),
			Audience:  []string{"https://api.oxyl.zhacked.me", "https://ingress.oxyl.zhacked.me"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "oxyl",
		},
	}), nil
}

type RefreshToken struct {
	Identifier    string
	Holder        *string
	AccessTokenId string

	Type JWTTokenType

	jwt.RegisteredClaims
}

func NewRefreshToken(identifier string, holder *string, accessTokenId string, tokenType JWTTokenType) (*RefreshToken, error) {
	if identifier == "" {
		return nil, errors.New("identifier is empty")
	}
	if accessTokenId == "" {
		return nil, errors.New("access token id is empty")
	}

	if tokenType != TokenTypeAgent && tokenType != TokenTypeUser {
		return nil, errors.New("token type is not valid")
	}

	if tokenType == TokenTypeAgent && (holder == nil || *holder == "") {
		return nil, errors.New("holder is nil")
	}

	return &RefreshToken{
		Identifier:    identifier,
		Holder:        holder,
		AccessTokenId: accessTokenId,
		Type:          tokenType,

		RegisteredClaims: jwt.RegisteredClaims{
			ID:        ulid.Make().String(),
			Audience:  []string{"https://api.oxyl.zhacked.me", "https://nexus.oxyl.zhacked.me"},
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 30)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "oxyl",
		},
	}, nil
}

type TokenPair struct {
	AccessToken struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	} `json:"access_token"`
	RefreshToken struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	} `json:"refresh_token"`
}
