package domain

// Account represents a financial account (checking, loan, card, etc.).
type Account struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	MonthlyFee float64 `json:"monthly_fee"`
	Type       string  `json:"type"`
}

// TariffAdjustmentRequest represents a tariff adjustment request for an account.
type TariffAdjustmentRequest struct {
	TransactionID string  `json:"transaction_id"`
	AccountID     string  `json:"account_id"`
	NewFee        float64 `json:"new_fee"`
	Status        string  `json:"status"`
}
