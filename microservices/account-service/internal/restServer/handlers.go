package restServer

import (
	"context"
	"net/http"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/dataTypes"
	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("k8ssecretkey")

func SignupHandler(c *gin.Context, mongoClient *mongodb.MongoDB) {
	log.Info().Msg("Registering new user.")

	var newUser dataTypes.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Error().Msg("Invalid request payload." + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload." + err.Error()})
		return
	}

	// Check if the user already exists by email or username
	if exists, err := checkUserExists(mongoClient, c.Request.Context(), &newUser); err != nil {
		log.Error().Msgf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	} else if exists {
		log.Error().Msg("User already registered.")
		c.JSON(http.StatusConflict, gin.H{"error": "User already registered."})
		return
	}

	newUser.Password = getHash([]byte(newUser.Password))
	_, err := mongoClient.Collections["USERS"].InsertOne(c.Request.Context(), newUser)
	if err != nil {
		log.Error().Msgf("Error inserting new user to database %v", err)
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully."})
	log.Info().Msg("User registered successfully.")
}

func LoginHandler(c *gin.Context, mongoClient *mongodb.MongoDB) {
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
	if err := findUserByEmail(mongoClient, c.Request.Context(), loginRequest.Email, &storedUser); err != nil {
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

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Error().Msgf("Error generating hash from password. %v", err)
	}
	return string(hash)
}

// Helper functions are going to be moved to different package to manage all all accounts operations.
func checkUserExists(mongoClient *mongodb.MongoDB, ctx context.Context, user *dataTypes.User) (bool, error) {
	filter := bson.M{"email": user.Email}
	var existingUser dataTypes.User
	err := mongoClient.Collections["USERS"].FindOne(ctx, filter).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments && err != mongo.ErrNilDocument {
		return false, err
	}
	return (err != mongo.ErrNoDocuments && err != mongo.ErrNilDocument), nil
}

func findUserByEmail(mongoClient *mongodb.MongoDB, ctx context.Context, email string, user *dataTypes.User) error {
	filter := bson.M{"email": email}
	result := mongoClient.Collections["USERS"].FindOne(ctx, filter)
	return result.Decode(user)
}

// checkPasswordHash verifies if the provided password matches the hashed password in the database
func checkPasswordHash(password, hashedPassword string) bool {
	// bcrypt.CompareHashAndPassword returns nil on a successful match
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func generateJWTToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Error().Msgf("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}
