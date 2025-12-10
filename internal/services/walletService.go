package services

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"whotterre/argent/internal/config"
	"whotterre/argent/internal/customErrors"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/models"
	"whotterre/argent/internal/repositories"
	"whotterre/argent/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService interface {
	DepositWallet(input dto.DepositWalletRequest, userID uuid.UUID) (*dto.DepositWalletResponse, error)
	GetBalance(userID uuid.UUID) (float64, error)
	Transfer(userID uuid.UUID, receiverWalletID string, amount float64) error
	GetTransactions(userID uuid.UUID) ([]models.Transaction, error)
	ProcessWebhook(payload []byte, signature string) error
	GetDepositStatus(reference string) (map[string]interface{}, error)
}

type walletService struct {
	walletRepo      repositories.WalletRepository
	transactionRepo repositories.TransactionRepository
	userRepo        repositories.UserRepository
	paystackSecret  string
	db              *gorm.DB
	config          config.Config
}

func NewWalletService(walletRepo repositories.WalletRepository, transactionRepo repositories.TransactionRepository, userRepo repositories.UserRepository, paystackSecret string, db *gorm.DB, cfg config.Config) WalletService {
	return &walletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		paystackSecret:  paystackSecret,
		db:              db,
		config:          cfg,
	}
}

func (s *walletService) DepositWallet(input dto.DepositWalletRequest, userID uuid.UUID) (*dto.DepositWalletResponse, error) {
	if input.Amount <= 0 {
		return nil, customErrors.ErrInsufficientFunds
	}

	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	// Generate reference
	ref := utils.GenRefString()

	// Create transaction
	transaction := &models.Transaction{
		ReceiverID: userID,
		Amount:     input.Amount,
		Type:       "deposit",
		Status:     "pending",
		Reference:  ref,
	}
	err = s.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Call Paystack
	payload := map[string]interface{}{
		"amount":       int(input.Amount * 100),
		"email":        user.Email,
		"reference":    ref,
		"callback_url": s.config.BaseURL + "/wallet/deposit/callback",
	}
	resp, err := s.callPaystack("transaction/initialize", payload)
	if err != nil {
		return nil, err
	}

	data := resp["data"].(map[string]interface{})

	return &dto.DepositWalletResponse{
		Reference:        ref,
		AuthorizationURL: data["authorization_url"].(string),
	}, nil
}

func (s *walletService) GetBalance(userID uuid.UUID) (float64, error) {
	wallet, err := s.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (s *walletService) Transfer(userID uuid.UUID, receiverWalletID string, amount float64) error {
	// Get sender wallet
	senderWallet, err := s.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		return err
	}
	if senderWallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	// Parse receiver ID
	receiverID, err := uuid.Parse(receiverWalletID)
	if err != nil {
		return errors.New("invalid receiver wallet ID")
	}

	// Get receiver wallet
	receiverWallet, err := s.walletRepo.GetWalletByUserID(receiverID)
	if err != nil {
		return err
	}

	// Atomic transfer
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Deduct from sender
	err = s.walletRepo.UpdateBalance(userID, senderWallet.Balance-amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add to receiver
	err = s.walletRepo.UpdateBalance(receiverID, receiverWallet.Balance+amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create transaction
	transaction := &models.Transaction{
		SenderID:   &userID,
		ReceiverID: receiverID,
		Amount:     amount,
		Type:       "transfer",
		Status:     "success",
	}
	err = s.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *walletService) GetTransactions(userID uuid.UUID) ([]models.Transaction, error) {
	return s.transactionRepo.GetUserTransactions(userID)
}

func (s *walletService) ProcessWebhook(payload []byte, signature string) error {
	// Validate signature
	expectedSignature := hmac.New(sha512.New, []byte(s.paystackSecret))
	expectedSignature.Write(payload)
	if !hmac.Equal([]byte(signature), expectedSignature.Sum(nil)) {
		return errors.New("invalid signature")
	}

	var data map[string]interface{}
	json.Unmarshal(payload, &data)

	event := data["event"].(string)
	if event != "charge.success" {
		return nil // ignore other events
	}

	reference := data["data"].(map[string]interface{})["reference"].(string)

	// Find transaction
	transaction, err := s.transactionRepo.GetTransactionByReference(reference)
	if err != nil {
		return err
	}

	if transaction.Status == "success" {
		return nil // idempotent
	}

	// Update status
	err = s.transactionRepo.UpdateTransactionStatus(transaction.ID, "success")
	if err != nil {
		return err
	}

	// Credit wallet
	wallet, err := s.walletRepo.GetWalletByUserID(transaction.ReceiverID)
	if err != nil {
		return err
	}
	newBalance := wallet.Balance + transaction.Amount
	err = s.walletRepo.UpdateBalance(transaction.ReceiverID, newBalance)
	if err != nil {
		return err
	}

	return nil
}

func (s *walletService) GetDepositStatus(reference string) (map[string]interface{}, error) {
	transaction, err := s.transactionRepo.GetTransactionByReference(reference)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"reference": reference,
		"status":    transaction.Status,
		"amount":    transaction.Amount,
	}, nil
}

func (s *walletService) callPaystack(endpoint string, payload map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.paystack.co/" + endpoint
	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(data)))
	req.Header.Set("Authorization", "Bearer "+s.paystackSecret)
	log.Println(s.paystackSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Paystack error: %v", result)
	}
	return result, nil
}
