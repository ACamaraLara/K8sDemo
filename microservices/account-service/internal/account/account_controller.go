package account

import (
	"account-service/internal/encryption"
	"account-service/internal/model"
	"errors"

	"context"
	"fmt"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/mongodb"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountController struct {
	mongoClient *mongodb.MongoDB
}

func NewAccountController(mongoClient *mongodb.MongoDB) *AccountController {
	return &AccountController{mongoClient: mongoClient}
}

func (s *AccountController) Signup(ctx context.Context, user *model.User) error {
	// Check if user already exists
	exists, err := s.checkUserExists(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("account already registered: %s", user.Email)
	}
	// Hash password
	user.Password, err = encryption.GetHash([]byte(user.Password))
	if err != nil {
		return err
	}

	_, err = s.mongoClient.Collections["USERS"].InsertOne(ctx, user)
	if err != nil {
		log.Error().Msgf("Error inserting new user to database %v", err)
		return err
	}
	return err
}

func (s *AccountController) Login(ctx context.Context, email, password string) (*model.User, error) {
	var storedUser model.User
	exist, err := s.findUser(ctx, email, &storedUser)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, errors.New("user not registered")
	}
	if !encryption.CheckPasswordHash(password, storedUser.Password) {
		return nil, errors.New("incorrect password")
	}
	return &storedUser, nil
}

func (s *AccountController) findUser(ctx context.Context, email string, user *model.User) (bool, error) {
	filter := bson.M{"email": email}
	result := s.mongoClient.Collections["USERS"].FindOne(ctx, filter)

	// If no document is found, return false and no error
	if err := result.Decode(user); err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// If user found, return true
	return true, nil
}

func (s *AccountController) checkUserExists(ctx context.Context, email string) (bool, error) {
	exists, err := s.findUser(ctx, email, &model.User{}) // We don't need to populate the user object
	if err != nil {
		return false, err
	}
	return exists, nil
}
