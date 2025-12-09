package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type APIKey struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Name        string         `gorm:"not null" json:"name"`
	HashedKey   string         `gorm:"not null" json:"-"`
	Permissions pq.StringArray `gorm:"type:text[];not null" json:"permissions"`
	ExpiresAt   time.Time      `gorm:"not null" json:"expires_at"`
	IsRevoked   bool           `gorm:"default:false" json:"is_revoked"`
	CreatedAt   time.Time      `gorm:"default:now()" json:"created_at"`
}

func (APIKey) TableName() string {
	return "api_keys"
}
