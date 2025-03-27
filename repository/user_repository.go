package repository

import (
	"mailcast/database"
	"mailcast/models"
)

// CreateUser inserts a new user into the database
func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

// GetUserByEmail finds a user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).Error
	return user, err
}
