package repositories

import (
	"log"
	"whotterre/argent/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWalletByUserID(userID uuid.UUID) (*models.Wallet, error)
	CreateWallet(wallet *models.Wallet) error
	UpdateBalance(walletID uuid.UUID, newBalance float64) error
	GetBalance(userID uuid.UUID) (float64, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) GetWalletByUserID(userID uuid.UUID) (*models.Wallet, error) {
	var wallet *models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		log.Println("Failed to get wallet by user ID:", err)
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepository) CreateWallet(wallet *models.Wallet) error {
	if err := r.db.Create(wallet).Error; err != nil {
		log.Println("Failed to create wallet:", err)
		return err
	}
	return nil
}

func (r *walletRepository) UpdateBalance(walletID uuid.UUID, newBalance float64) error {
	if err := r.db.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", newBalance).Error; err != nil {
		log.Println("Failed to update balance:", err)
		return err
	}
	return nil
}

func (r *walletRepository) GetBalance(userID uuid.UUID) (float64, error) {
	var balance float64
	if err := r.db.Model(&models.Wallet{}).Where("user_id = ?", userID).Select("balance").Scan(&balance).Error; err != nil {
		log.Println("Failed to get balance:", err)
		return 0, err
	}
	return balance, nil
}
