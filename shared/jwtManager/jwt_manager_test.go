package jwtManager

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr string
	}{
		{
			name: "valid config",
			config: &Config{
				SecretKey:         "test-secret",
				AccessTokenExpiry: 1 * time.Hour,
				Issuer:            "test-issuer",
			},
			wantErr: "",
		},
		{
			name: "missing secret key",
			config: &Config{
				SecretKey:         "",
				AccessTokenExpiry: 1 * time.Hour,
				Issuer:            "test-issuer",
			},
			wantErr: "secret key must be provided",
		},
		{
			name: "missing issuer",
			config: &Config{
				SecretKey:         "test-secret",
				AccessTokenExpiry: 1 * time.Hour,
				Issuer:            "",
			},
			wantErr: "jwt issuer must be provided",
		},
		{
			name: "invalid expiration time",
			config: &Config{
				SecretKey:         "test-secret",
				AccessTokenExpiry: 0,
				Issuer:            "test-issuer",
			},
			wantErr: "expiration time should be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.config)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, manager)
			} else {
				assert.EqualError(t, err, tt.wantErr)
				assert.Nil(t, manager)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	config := &Config{
		SecretKey:         "test-secret",
		AccessTokenExpiry: 1 * time.Hour,
		Issuer:            "test-issuer",
	}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	token, err := manager.GenerateToken("123", "test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Decode the token and validate the claims
	parsedToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(*UserClaims)
	assert.True(t, ok)
	assert.Equal(t, "123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.WithinDuration(t, time.Now().Add(1*time.Hour), claims.ExpiresAt.Time, 2*time.Second)
}

func TestValidateToken(t *testing.T) {
	config := &Config{
		SecretKey:         "test-secret",
		AccessTokenExpiry: 1 * time.Hour,
		Issuer:            "test-issuer",
	}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	token, err := manager.GenerateToken("123", "test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	claims, err := manager.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "test-issuer", claims.Issuer)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	config := &Config{
		SecretKey:         "test-secret",
		AccessTokenExpiry: 1 * time.Hour,
		Issuer:            "test-issuer",
	}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	// Validate an invalid token
	invalidToken := "invalid.token.value"
	claims, err := manager.ValidateToken(invalidToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	config := &Config{
		SecretKey:         "test-secret",
		AccessTokenExpiry: 1 * time.Second, // Expired token
		Issuer:            "test-issuer",
	}
	manager, err := NewManager(config)
	assert.NoError(t, err)

	token, err := manager.GenerateToken("123", "test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	time.Sleep(2 * time.Second)
	// Validate the expired token
	claims, err := manager.ValidateToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
	assert.Nil(t, claims)
}
