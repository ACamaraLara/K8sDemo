package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHash(t *testing.T) {
	password := "password123"
	hash, err := GetHash([]byte(password))

	// Ensure the hash is not empty and not equal to the plain password
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Test for empty password
	emptyHash, err := GetHash([]byte(""))
	assert.NoError(t, err)
	assert.NotEmpty(t, emptyHash)
	assert.NotEqual(t, "", emptyHash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	hashedPassword, _ := GetHash([]byte(password))

	// Test correct password
	t.Run("Correct password", func(t *testing.T) {
		isValid := CheckPasswordHash(password, hashedPassword)
		assert.True(t, isValid, "Password should be valid")
	})

	// Test incorrect password
	t.Run("Incorrect password", func(t *testing.T) {
		incorrectPassword := "wrongpassword"
		isValid := CheckPasswordHash(incorrectPassword, hashedPassword)
		assert.False(t, isValid, "Password should not be valid")
	})
}
