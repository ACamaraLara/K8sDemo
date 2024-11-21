package mongodb

import (
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAddFlagsParams(t *testing.T) {
	// Setup environment variables for testing
	os.Setenv("MONGODB_HOST", "testhost")
	os.Setenv("MONGODB_PORT", "27018")
	os.Setenv("MONGODB_DATABASE", "testdb")
	os.Setenv("MONGODB_USER", "testuser")
	os.Setenv("MONGODB_PASSWD", "testpassword")
	os.Setenv("MONGODB_COLLECTIONS", "col1,col2")

	// Create MongoConfig instance
	cfg := &MongoConfig{}

	// Call the function to set flags and load environment variables
	cfg.AddFlagsParams()

	// Ensure the fields are correctly set from environment variables
	assert.Equal(t, "testhost", cfg.Host)
	assert.Equal(t, "27018", cfg.Port)
	assert.Equal(t, "testdb", cfg.DbName)
	assert.Equal(t, "testuser", cfg.User)
	assert.Equal(t, "testpassword", cfg.Passwd)
	assert.Equal(t, []string{"col1", "col2"}, cfg.Collections)

	// Cleanup environment variables
	os.Unsetenv("MONGODB_HOST")
	os.Unsetenv("MONGODB_PORT")
	os.Unsetenv("MONGODB_DATABASE")
	os.Unsetenv("MONGODB_USER")
	os.Unsetenv("MONGODB_PASSWD")
	os.Unsetenv("MONGODB_COLLECTIONS")
}

func TestGetURL(t *testing.T) {
	// Define a struct to hold test case data
	type testCase struct {
		name        string
		config      *MongoConfig
		expectedURL string
	}

	// Define the test cases
	testCases := []testCase{
		{
			name: "WithoutUserAndPasswd",
			config: &MongoConfig{
				Host:   DefaultMongoDBHost,
				Port:   DefaultMongoDBPort,
				DbName: DefaultMongoDBName,
			},
			expectedURL: "mongodb://" + DefaultMongoDBHost + ":" + DefaultMongoDBPort + "/" + DefaultMongoDBName,
		},
		{
			name: "WithUserAndPasswd",
			config: &MongoConfig{
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
