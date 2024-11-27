package restServer

import (
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/dataTypes"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/database"
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("k8ssecretkey")

func SignupHandler(c *gin.Context, dbClient *database.DBManager) {
	log.Info().Msg("Registering new user.")

	var newUser dataTypes.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Error().Msg("Invalid request payload." + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload." + err.Error()})
		return
	}

	// Check if the user already exists by email or username
	if err := dbClient.FindOne(c.Request.Context(), "USERS", nil,
		map[string]interface{}{"email": newUser.Email}); err == nil {
		log.Error().Msg("User already registered.")
		c.JSON(http.StatusConflict, gin.H{"error": "User already registered."})
		return
	}

	newUser.Password = getHash([]byte(newUser.Password))
	if err := dbClient.InsertOne(c.Request.Context(), "USERS", newUser); err != nil {
		log.Error().Msgf("Error inserting new user to database %v", err)
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully."})
	log.Info().Msg("User registered successfully.")
}

func LoginHandler(c *gin.Context, dbClient *database.DBManager) {
	log.Info().Msg("User login attempt.")

	var loginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Error().Msg("Invalid login request payload.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload."})
		return
	}

	var storedUser dataTypes.User
	if err := dbClient.FindOne(c.Request.Context(), "USERS", &storedUser,
		map[string]interface{}{"email": loginRequest.Email}); err != nil {
		log.Error().Msgf("User not found: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password."})
		return
	}

	// Verify password
	if !checkPasswordHash(loginRequest.Password, storedUser.Password) {
		log.Error().Msg("Incorrect password.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect email or password."})
		return
	}

	// Generate JWT token on successful authentication
	token, err := generateJWTToken()
	if err != nil {
		log.Error().Msgf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token."})
		return
	}

	log.Info().Msg("User logged in successfully.")
	c.JSON(http.StatusOK, gin.H{"authToken": token})
}

// :ToDo: Create package to manage password encryption.

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Error().Msgf("Error generating hash from password. %v", err)
	}
	return string(hash)
}

// checkPasswordHash verifies if the provided password matches the hashed password in the database
func checkPasswordHash(password, hashedPassword string) bool {
	// bcrypt.CompareHashAndPassword returns nil on a successful match
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// :ToDo: Create package to manage jwt token.

func generateJWTToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Error().Msgf("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}
