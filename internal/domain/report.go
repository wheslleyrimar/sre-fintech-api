package domain

// Report holds fintech summary metrics: total accounts, counts by type, top accounts by fee.
type Report struct {
	TotalAccounts  int                `json:"total_accounts"`
	TotalsByType   map[string]int     `json:"totals_by_type"`
	Top100ByFee    []Account          `json:"top_100_by_fee"`
}
