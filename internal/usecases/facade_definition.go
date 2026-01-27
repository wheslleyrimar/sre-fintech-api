package usecases

import (
	"context"
	"sre/internal/domain"
)

// AccountSearcher searches accounts by term.
type AccountSearcher interface {
	SearchByTerm(ctx context.Context, term string) ([]domain.Account, error)
}

// AdjustmentFlowProcessor starts the tariff adjustment approval flow.
type AdjustmentFlowProcessor interface {
	BeginFlow(ctx context.Context, input domain.TariffAdjustmentRequest, callbackURL string) error
}
