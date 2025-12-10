package handlers

import (
	"net/http"
	"whotterre/argent/internal/dto"
	"whotterre/argent/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService services.WalletService
}

func NewWalletHandler(walletService services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	var req dto.DepositWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	response, err := h.walletService.DepositWallet(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	balance, err := h.walletService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.BalanceResponse{Balance: balance})
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	var req dto.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	err := h.walletService.Transfer(userID, req.WalletNumber, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TransferResponse{Status: "success", Message: "Transfer completed"})
}

func (h *WalletHandler) GetTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	transactions, err := h.walletService.GetTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.TransactionResponse
	for _, t := range transactions {
		response = append(response, dto.TransactionResponse{
			Type:   t.Type,
			Amount: t.Amount,
			Status: t.Status,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *WalletHandler) GetDepositStatus(c *gin.Context) {
	reference := c.Param("reference")

	status, err := h.walletService.GetDepositStatus(reference)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *WalletHandler) Webhook(c *gin.Context) {
	payload, _ := c.GetRawData()
	signature := c.GetHeader("x-paystack-signature")

	err := h.walletService.ProcessWebhook(payload, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (h *WalletHandler) DepositCallback(c *gin.Context) {
	reference := c.Query("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reference missing"})
		return
	}

	status, err := h.walletService.GetDepositStatus(reference)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}
