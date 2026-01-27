package usecases

import (
	"context"
	"sre/internal/domain"
)

// AccountRepository reads and updates accounts in the fintech API.
type AccountRepository interface {
	UpdateFee(ctx context.Context, acc domain.Account, newFee float64) error
	Get(ctx context.Context, id domain.Account) (domain.Account, error)
}

// TariffAdjustmentRepository manages tariff adjustment requests.
type TariffAdjustmentRepository interface {
	Create(ctx context.Context, input domain.TariffAdjustmentRequest, callbackURL string) error
	GetLastByAccount(ctx context.Context, acc domain.Account) (*domain.TariffAdjustmentRequest, error)
	AllByAccount(ctx context.Context, acc domain.Account) ([]domain.TariffAdjustmentRequest, error)
}
