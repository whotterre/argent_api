package dto

type DepositWalletRequest struct {
	Amount float64 `json:"amount"`
}

type DepositWalletResponse struct {
	Reference        string `json:"reference"`
	AuthorizationURL string `json:"authorization_url"`
}

type TransferRequest struct {
	WalletNumber string  `json:"wallet_number"`
	Amount       float64 `json:"amount"`
}

type TransferResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

type TransactionResponse struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type DepositStatusResponse struct {
	Reference string  `json:"reference"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
}
