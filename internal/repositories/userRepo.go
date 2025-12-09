package repositories

import (
	"log"
	"github.com/google/uuid"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindOrCreateUser(newUser *dto.CreateNewUserRequest) (*models.User, error)
	GetUserById(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByGoogleID(googleID string) (*models.User, error)
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
	var user models.User
	// Check if user exists
	err := r.db.Preload("Wallet").Where("email = ?", newUser.Email).First(&user).Error
	if err == nil {
		return &user, nil
	}

	// If not found, create
	tx := r.db.Begin()
	user = models.User{
		GoogleID:  newUser.GoogleID,
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to create user:", err)
		return nil, err
	}

	// Create Wallet
	wallet := models.Wallet{
		UserID:  user.ID,
		Balance: 0,
	}
	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to create wallet:", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	user.Wallet = wallet
	return &user, nil
}

func (r *userRepository) GetUserById(id uuid.UUID) (*models.User, error){
	var user *models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Println("Failed to get user by ID")
		return nil, err
	}
	return user, nil
}
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		log.Println("Failed to get user by email")
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByGoogleID(googleID string) (*models.User, error) {
	var user *models.User
	if err := r.db.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		log.Println("Failed to get user by Google ID")
		return nil, err
	}
	return user, nil
}
