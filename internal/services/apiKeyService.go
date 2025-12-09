package services

import (
	"errors"
	"log"
	"slices"
	"time"
	"whotterre/argent/internal/customErrors"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/models"
	"whotterre/argent/internal/repositories"
	"whotterre/argent/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type APIKeyService interface {
	CreateAPIKey(input dto.CreateAPIKeyRequest, userID uuid.UUID) (*dto.CreateAPIKeyResponse, error)
	ValidateAPIKey(apiKey string, userID uuid.UUID, requiredPermission string) (*models.APIKey, error)
	RolloverAPIKey(input *dto.RolloverAPIKeyRequest) (*dto.RolloverAPIKeyResponse, error)
}

type apiKeyService struct {
	apiKeyRepo repositories.APIKeyRepository
}

func NewAPIKeyService(apiKeyRepo repositories.APIKeyRepository) APIKeyService {
	return &apiKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

func (s *apiKeyService) CreateAPIKey(input dto.CreateAPIKeyRequest, userID uuid.UUID) (*dto.CreateAPIKeyResponse, error) {
	// Ensure user has < 5 active API Keys
	userAPIKeys, err := s.apiKeyRepo.GetActiveAPIKeysByUserID(userID)
	if err != nil {
		log.Println("Failed to get user's active API keys ", err.Error())
		return nil, err
	}

	if len(userAPIKeys) == 5 {
		return nil, customErrors.ErrorActiveAPIKeysExceeded
	}

	apiKey := utils.GenerateNewAPIKeyString()
	expiryDate, err := utils.ExpiryStringToTimestamp(input.Expiry)
	if err != nil {
		return nil, err
	}

	// Validate permissions
	validPermissions := []string{"deposit", "transfer", "read"}
	for _, perm := range input.Permissions {
		if !slices.Contains(validPermissions, perm) {
			return nil, customErrors.ErrInvalidPermission
		}
	}
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to generate hash for API key:", err)
		return nil, err
	}

	newAPIKey := models.APIKey{
		UserID:      userID,
		Name:        input.Name,
		HashedKey:   string(hashedKey),
		Permissions: input.Permissions,
		ExpiresAt:   expiryDate,
	}

	err = s.apiKeyRepo.CreateAPIKey(&newAPIKey)
	if err != nil {
		log.Println("Failed to create API key in DB:", err)
		return nil, err
	}

	result := dto.CreateAPIKeyResponse{
		APIKey:    apiKey,
		ExpiresAt: expiryDate,
	}

	return &result, nil
}

func (s *apiKeyService) ValidateAPIKey(apiKey string, userID uuid.UUID, requiredPermission string) (*models.APIKey, error) {
	userAPIKeys, err := s.apiKeyRepo.GetActiveAPIKeysByUserID(userID)
	if err != nil {
		log.Println("Failed to get user's active API keys:", err)
		return nil, err
	}

	for _, key := range userAPIKeys {
		if err := bcrypt.CompareHashAndPassword([]byte(key.HashedKey), []byte(apiKey)); err == nil {
			// Key matches, check permission
			if !containsPermission(key.Permissions, requiredPermission) {
				return nil, errors.New("insufficient permissions")
			}
			return &key, nil
		}
	}

	return nil, errors.New("invalid API key")
}

func (s *apiKeyService) RolloverAPIKey(input *dto.RolloverAPIKeyRequest) (*dto.RolloverAPIKeyResponse, error) {
	allAPIKeys, err := s.apiKeyRepo.GetAllNonRevokedAPIKeys()
	if err != nil {
		log.Println("Failed to get all active API keys:", err)
		return nil, err
	}

	var expiredKey *models.APIKey
	for _, key := range allAPIKeys {
		if err := bcrypt.CompareHashAndPassword([]byte(key.HashedKey), []byte(input.ExpiredAPIKey)); err == nil {
			expiredKey = &key
			break
		}
	}

	if expiredKey == nil {
		return nil, customErrors.ErrNonExistentAPIKey
	}

	if !time.Now().After(expiredKey.ExpiresAt) {
		return nil, customErrors.ErrRollingOverNotExpiredKey
	}

	// Revoke the old key
	err = s.apiKeyRepo.RevokeAPIKey(expiredKey.ID)
	if err != nil {
		return nil, err
	}

	newAPIKey := dto.CreateAPIKeyRequest{
		Name:        expiredKey.Name,
		Permissions: expiredKey.Permissions,
		Expiry:      input.Expiry,
	}

	createdKey, err := s.CreateAPIKey(newAPIKey, expiredKey.UserID)
	if err != nil {
		return nil, err
	}

	response := dto.RolloverAPIKeyResponse{
		APIKey:    createdKey.APIKey,
		ExpiresAt: createdKey.ExpiresAt,
	}
	return &response, nil
}

func containsPermission(permissions []string, permission string) bool {
	return slices.Contains(permissions, permission)
}
