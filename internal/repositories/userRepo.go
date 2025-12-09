package repositories

import (
	"log"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindOrCreateUser(newUser *dto.CreateNewUserRequest) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindOrCreateUser(newUser *dto.CreateNewUserRequest) (*models.User, error) {
	existingUser, err := r.GetUserByEmail(newUser.Email)
	if err == nil {
		return existingUser, nil
	}

	if err != gorm.ErrRecordNotFound {
		log.Println("Error checking for existing user:", err)
		return nil, err
	}

	user := models.User{
		GoogleID: newUser.GoogleID,
		Email:    newUser.Email,
		FullName: newUser.FullName,
		IsActive: true,
	}

	if err := r.db.Create(&user).Error; err != nil {
		log.Println("Failed to create user:", err)
		return nil, err
	}

	return &user, nil
}
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		log.Println("Failed to get user by email")
		return nil, err
	}
	return user, nil
}
