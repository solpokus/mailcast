package services

import (
	"errors"
	"fmt"
	"mailcast/models"
	"mailcast/repository"
)

// RegisterUser handles user creation logic
func RegisterUser(name, email string) (models.User, error) {
	// Check if user already exists
	existingUser, _ := repository.GetUserByEmail(email)
	if existingUser.ID != 0 {
		return models.User{}, errors.New("email already registered")
	}

	// Create a new user
	user := models.User{FirstName: name, Email: email}
	err := repository.CreateUser(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserByEmailService fetches a user by email (wrapper over repository)
func GetUserByEmailService(email string) (models.User, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}
	fmt.Println("user firstname", user.FirstName)
	return user, nil
}
