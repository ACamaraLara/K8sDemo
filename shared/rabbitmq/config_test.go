package rabbitmq

import (
	"flag"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func cleanupFlags() {
	// Reset the flag set to avoid conflicts between tests
	flag.CommandLine = flag.NewFlagSet(flag.CommandLine.Name(), flag.PanicOnError)
}

func TestRabbitConfig_WithFlags(t *testing.T) {
	cfg := &RabbitConfig{}
	cfg.AddFlagsParams()

	// Set flags to override the environment variables and default values
	flag.Set("rabbit-host", "flag-host")
	flag.Set("rabbit-port", "1234")
	flag.Set("rabbit-user", "flag-user")
	flag.Set("rabbit-passwd", "flag-pass")

	flag.Parse()

	assert.Equal(t, "flag-host", cfg.Host)
	assert.Equal(t, "1234", cfg.Port)
	assert.Equal(t, "flag-user", cfg.User)
	assert.Equal(t, "flag-pass", cfg.Passwd)

	t.Cleanup(cleanupFlags)
}

func TestRabbitConfig_DefaultValues(t *testing.T) {
	// Reset the environment variables to ensure defaults are used
	os.Clearenv()

	cfg := &RabbitConfig{}
	cfg.AddFlagsParams()

	assert.Equal(t, DefaultRabbitHost, cfg.Host)
	assert.Equal(t, DefaultRabbitPort, cfg.Port)
	assert.Equal(t, DefaultRabbitUser, cfg.User)
	assert.Equal(t, DefaultRabbitPass, cfg.Passwd)

	t.Cleanup(cleanupFlags)
}

func TestRabbitConfig_WithEnvironmentVariables(t *testing.T) {
	os.Setenv("RABBITMQ_HOST", "host")
	os.Setenv("RABBITMQ__PORT", "1234")
	os.Setenv("USERNAME", "user")
	os.Setenv("PASSWORD", "pass")

	cfg := &RabbitConfig{}
	cfg.AddFlagsParams()

	assert.Equal(t, "host", cfg.Host)
	assert.Equal(t, "1234", cfg.Port)
	assert.Equal(t, "user", cfg.User)
	assert.Equal(t, "pass", cfg.Passwd)

}

func TestRabbitConfig_GetURL(t *testing.T) {
	testCases := []struct {
		name     string
		cfg      RabbitConfig
		expected string
	}{
		{
			name: "WithCredentials",
			cfg: RabbitConfig{
				Host:   "localhost",
				Port:   "1234",
				User:   "user",
				Passwd: "password",
			},
			expected: "amqp://user:password@localhost:1234/",
		},
		{
			name: "WithoutCredentials",
			cfg: RabbitConfig{
				Host: "localhost",
				Port: "1234",
			},
			expected: "amqp://localhost:1234/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := tc.cfg.GetURL()
			assert.Equal(t, tc.expected, url)
		})
	}
}
