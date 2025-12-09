package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SenderID   *uuid.UUID `gorm:"type:uuid" json:"sender_id"` // Pointer to allow null for deposits (system -> user)
	Sender     *User      `gorm:"foreignKey:SenderID;references:ID" json:"sender"`
	ReceiverID uuid.UUID  `gorm:"type:uuid;not null" json:"receiver_id"`
	Receiver   User       `gorm:"foreignKey:ReceiverID;references:ID" json:"receiver"`
	Amount     float64    `gorm:"not null" json:"amount"`
	Type       string     `gorm:"not null" json:"type"`    // 'deposit', 'transfer'
	Status     string     `gorm:"not null" json:"status"`  // "success|failed|pending"
	Reference  string     `gorm:"unique" json:"reference"` // Paystack reference
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}
