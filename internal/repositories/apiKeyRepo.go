package repositories

import (
	"log"
	"time"
	"whotterre/argent/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIKeyRepository interface {
	CreateAPIKey(apiKey *models.APIKey) error
	GetAllNonRevokedAPIKeys() ([]models.APIKey, error)
	GetActiveAPIKeysByUserID(userID uuid.UUID) ([]models.APIKey, error)
	GetAPIKeyByID(id uuid.UUID) (*models.APIKey, error)
	RevokeAPIKey(id uuid.UUID) error
	GetExpiredKeyByID(id uuid.UUID) (*models.APIKey, error)
}

type apiKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{
		db: db,
	}
}

func (r *apiKeyRepository) CreateAPIKey(apiKey *models.APIKey) error {
	if err := r.db.Create(apiKey).Error; err != nil {
		log.Println("Failed to create API key:", err)
		return err
	}
	return nil
}

func (r *apiKeyRepository) GetAllNonRevokedAPIKeys() ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	if err := r.db.Where("is_revoked = false").
		Find(&apiKeys).Error; err != nil {
		log.Println("Failed to get all non-revoked API keys:", err)
		return nil, err
	}
	return apiKeys, nil
}

func (r *apiKeyRepository) GetActiveAPIKeysByUserID(userID uuid.UUID) ([]models.APIKey, error) {
	var apiKeys []models.APIKey
	if err := r.db.Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Find(&apiKeys).Error; err != nil {
		log.Println("Failed to get active API keys by user ID:", err)
		return nil, err
	}
	return apiKeys, nil
}

func (r *apiKeyRepository) GetAPIKeyByID(id uuid.UUID) (*models.APIKey, error) {
	var apiKey *models.APIKey
	if err := r.db.Where("id = ?", id).First(&apiKey).Error; err != nil {
		log.Println("Failed to get API key by ID:", err)
		return nil, err
	}
	return apiKey, nil
}

func (r *apiKeyRepository) RevokeAPIKey(id uuid.UUID) error {
	if err := r.db.Model(&models.APIKey{}).Where("id = ?", id).Update("is_revoked", true).Error; err != nil {
		log.Println("Failed to revoke API key:", err)
		return err
	}
	return nil
}

func (r *apiKeyRepository) GetExpiredKeyByID(id uuid.UUID) (*models.APIKey, error) {
	var apiKey *models.APIKey
	if err := r.db.Where("id = ? AND expires_at <= ?", id, time.Now()).First(&apiKey).Error; err != nil {
		log.Println("Failed to get expired key:", err)
		return nil, err
	}
	return apiKey, nil
}
