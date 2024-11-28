package jwtManager

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	conf *Config
}

// UserClaims represents custom claims that refers to the user that will request the new token.
type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewManager(config *Config) (*Manager, error) {
	log.Info().Msgf("Creating JWT manager.")
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	return &Manager{
		conf: config,
	}, nil
}

func validateConfig(config *Config) error {
	if config.AccessTokenExpiry <= 0 {
		return errors.New("expiration time should be greater than 0")
	} else if config.SecretKey == "" {
		return errors.New("secret key must be provided")
	} else if config.Issuer == "" {
		return errors.New("jwt issuer must be provided")
	}
	return nil
}

func (m *Manager) GenerateToken(userID, email string) (string, error) {
	log.Info().Msgf("Generating Token for user %s with email %s", userID, email)
	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.conf.AccessTokenExpiry)),
			Issuer:    m.conf.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.conf.SecretKey))
}

func (m *Manager) ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(m.conf.SecretKey), nil
	})
	// Checks if token is expired or invalid and return's error.
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
