package config

import (
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAddFlagsParams(t *testing.T) {
	// Setup environment variables for testing
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "27018")
	os.Setenv("DB_DATABASE", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWD", "testpassword")
	os.Setenv("DB_TABLES", "col1,col2")

	cfg := &DBConfig{}

	cfg.AddFlagsParams()

	assert.Equal(t, "testhost", cfg.Host)
	assert.Equal(t, "27018", cfg.Port)
	assert.Equal(t, "testdb", cfg.DbName)
	assert.Equal(t, "testuser", cfg.User)
	assert.Equal(t, "testpassword", cfg.Passwd)
	assert.Equal(t, []string{"col1", "col2"}, cfg.Collections)

	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_DATABASE")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWD")
	os.Unsetenv("DB_TABLES")
}

func TestGetURL(t *testing.T) {
	type testCase struct {
		name        string
		config      *DBConfig
		expectedURL string
	}

	// Define the test cases
	testCases := []testCase{
		{
			name: "WithoutUserAndPasswd",
			config: &DBConfig{
				Host:   DefaultMongoDBHost,
				Port:   DefaultMongoDBPort,
				DbName: DefaultMongoDBName,
			},
			expectedURL: "mongodb://" + DefaultMongoDBHost + ":" + DefaultMongoDBPort + "/" + DefaultMongoDBName,
		},
		{
			name: "WithUserAndPasswd",
			config: &DBConfig{
				Host:   DefaultMongoDBHost,
				Port:   DefaultMongoDBPort,
				DbName: DefaultMongoDBName,
				User:   "testuser",
				Passwd: "testpassword",
			},
			expectedURL: "mongodb://testuser:testpassword@" + DefaultMongoDBHost + ":" + DefaultMongoDBPort + "/" + DefaultMongoDBName,
		},
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetURL function
			url := tc.config.GetURL()

			// Check the generated URL
			assert.Equal(t, tc.expectedURL, url)
		})
	}
}
