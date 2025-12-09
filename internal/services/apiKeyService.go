package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"slices"
	"time"
	"whotterre/argent/internal/customErrors"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/models"
	"whotterre/argent/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type APIKeyService interface {
	CreateAPIKey(input dto.CreateAPIKeyRequest, userID uuid.UUID) (*dto.CreateAPIKeyResponse, error)
	ValidateAPIKey(apiKey string, requiredPermission string) (*models.APIKey, error)
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

	apiKey := generateNewAPIKeyString()
	expiryDate, err := expiryStringToTimestamp(input.Expiry)
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

func (s *apiKeyService) ValidateAPIKey(apiKey string, requiredPermission string) (*models.APIKey, error) {
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash API key for validation:", err)
		return nil, err
	}

	apiKeyRecord, err := s.apiKeyRepo.GetAPIKeyByHashedKey(string(hashedKey))
	if err != nil {
		log.Println("API key not found or invalid:", err)
		return nil, err
	}

	// Check if revoked or expired
	if apiKeyRecord.IsRevoked || time.Now().After(apiKeyRecord.ExpiresAt) {
		return nil, errors.New("API key is revoked or expired")
	}

	// Check permission
	if !containsPermission(apiKeyRecord.Permissions, requiredPermission) {
		return nil, errors.New("insufficient permissions")
	}

	return apiKeyRecord, nil
}

func containsPermission(permissions []string, permission string) bool {
	return slices.Contains(permissions, permission)
}

func generateNewAPIKeyString() string {
	prefix := "sk_live_"

	bytes := make([]byte, 45)
	rand.Read(bytes)
	end := base64.RawURLEncoding.EncodeToString(bytes)

	return prefix + end
}

func expiryStringToTimestamp(expiryStr string) (time.Time, error) {
	now := time.Now()
	switch expiryStr {
	case "1H":
		return now.Add(time.Hour), nil
	case "1D":
		return now.AddDate(0, 0, 1), nil 
	case "1M":
		return now.AddDate(0, 1, 0), nil 
	case "1Y":
		return now.AddDate(1, 0, 0), nil 
	default:
		return time.Time{}, errors.New("invalid expiry string") 
	}
}
