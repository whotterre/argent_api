package repositories

import (
	"log"
	"whotterre/argent/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetUserTransactions(userID uuid.UUID) ([]models.Transaction, error)
	GetTransactionByID(id uuid.UUID) (*models.Transaction, error)
	GetTransactionByReference(reference string) (*models.Transaction, error)
	UpdateTransactionStatus(id uuid.UUID, status string) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	if err := r.db.Create(transaction).Error; err != nil {
		log.Println("Failed to create transaction:", err)
		return err
	}
	return nil
}

func (r *transactionRepository) GetUserTransactions(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.Where("receiver_id = ? OR sender_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		log.Println("Failed to get user transactions:", err)
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction *models.Transaction
	if err := r.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		log.Println("Failed to get transaction by ID:", err)
		return nil, err
	}
	return transaction, nil
}

func (r *transactionRepository) GetTransactionByReference(reference string) (*models.Transaction, error) {
	var transaction *models.Transaction
	if err := r.db.Where("reference = ?", reference).First(&transaction).Error; err != nil {
		log.Println("Failed to get transaction by reference:", err)
		return nil, err
	}
	return transaction, nil
}

func (r *transactionRepository) UpdateTransactionStatus(id uuid.UUID, status string) error {
	if err := r.db.Model(&models.Transaction{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		log.Println("Failed to update transaction status:", err)
		return err
	}
	return nil
}
