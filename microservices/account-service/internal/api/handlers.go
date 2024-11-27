package api

import (
	"net/http"

	"account-service/internal/account"
	"account-service/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SignupHandler(c *gin.Context, accountService *account.AccountController) {
	var newUser model.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Error().Msgf("Invalid signup request payload %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signup request payload." + err.Error()})
		return
	}
	log.Info().Msgf("Registering New user %s", newUser.Email)

	// Use account service to handle the signup logic
	if err := accountService.Signup(c, &newUser); err != nil {
		log.Error().Msgf("Error signing up user %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Msgf("User registered successfully.")
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully."})
}

func LoginHandler(c *gin.Context, accountService *account.AccountController) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Error().Msgf("Invalid login request payload %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login request payload."})
		return
	}

	log.Info().Msgf("Logging user %s", loginRequest.Email)

	if err := accountService.Login(c, loginRequest.Email, loginRequest.Password); err != nil {
		log.Error().Msgf("Error logging in user %+v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Info().Msgf("User loged in successfully.")
	c.JSON(http.StatusOK, gin.H{"message": "User logged successfully!"})
}
