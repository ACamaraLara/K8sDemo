package jwtManager_test

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/jwtManager"
	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestConfig_AddFlagsParams(t *testing.T) {
	t.Run("missing JWT_SECRET_KEY should set an empty byte array", func(t *testing.T) {
		resetFlags()
		os.Clearenv()

		cfg := &jwtManager.Config{}
		cfg.AddFlagsParams()

		assert.Equal(t, "", cfg.SecretKey, "secretKey should be set to an empty byte array when JWT_SECRET_KEY is missing")
	})

	t.Run("valid JWT_SECRET_KEY environment variable", func(t *testing.T) {
		resetFlags()
		os.Clearenv()
		os.Setenv("JWT_SECRET_KEY", "valid-secret")

		cfg := &jwtManager.Config{}
		cfg.AddFlagsParams()

		assert.Equal(t, "valid-secret", cfg.SecretKey, "secretKey should match the JWT_SECRET_KEY environment variable")
	})

	t.Run("valid JWT_SECRET_KEY with custom flags and environment variables", func(t *testing.T) {
		resetFlags()
		os.Clearenv()
		os.Setenv("JWT_SECRET_KEY", "custom-secret")
		os.Setenv("JWT_EXPIRATION_HOURS", "6")
		os.Setenv("JWT_ISSUER", "custom-issuer")

		// Simulate command-line arguments
		os.Args = []string{
			"test", // Dummy program name
			"-jwt-issuer=flag-issuer",
		}

		cfg := &jwtManager.Config{}
		cfg.AddFlagsParams()
		flag.Parse()

		assert.Equal(t, 6*time.Hour, cfg.AccessTokenExpiry, "Command-line flags should override environment variables for jwt-expiration")
		assert.Equal(t, "flag-issuer", cfg.Issuer, "Command-line flags should override environment variables for jwt-issuer")
		assert.Equal(t, "custom-secret", cfg.SecretKey, "secretKey should match the JWT_SECRET_KEY environment variable")
	})
}
