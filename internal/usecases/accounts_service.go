package usecases

import (
	"context"

	"sre/internal/domain"

	"log/slog"
)

var _ AccountService = (*AccountServiceImpl)(nil)

// AccountService exposes fintech account and tariff-adjustment operations.
type AccountService interface {
	UpdateFee(ctx context.Context, accountID string) error
	SendTariffAdjustmentRequest(ctx context.Context, input domain.TariffAdjustmentRequest) error
	GetTariffAdjustments(ctx context.Context, acc domain.Account) ([]domain.TariffAdjustmentRequest, error)
	GetAccount(ctx context.Context, accountID string) (domain.Account, error)
}

// NewAccountService creates an AccountService.
func NewAccountService(
	accountRepo AccountRepository,
	adjustmentRepo TariffAdjustmentRepository,
	flowProcessor AdjustmentFlowProcessor,
	callbackBaseURL string,
) AccountService {
	return &AccountServiceImpl{
		accountRepo:   accountRepo,
		adjustmentRepo: adjustmentRepo,
		flowProcessor: flowProcessor,
		callbackURL:   callbackBaseURL + "/accounts/notifications",
	}
}

type AccountServiceImpl struct {
	accountRepo    AccountRepository
	adjustmentRepo TariffAdjustmentRepository
	flowProcessor  AdjustmentFlowProcessor
	callbackURL    string
}

func (s *AccountServiceImpl) SendTariffAdjustmentRequest(ctx context.Context, input domain.TariffAdjustmentRequest) error {
	slog.InfoContext(ctx, "sending tariff adjustment request", "input", input)
	go func() {
		if err := s.adjustmentRepo.Create(context.Background(), input, s.callbackURL); err != nil {
			slog.ErrorContext(ctx, "create tariff adjustment failed", "err", err)
		}
	}()
	go func() {
		if err := s.flowProcessor.BeginFlow(context.Background(), input, s.callbackURL); err != nil {
			slog.ErrorContext(ctx, "begin adjustment flow failed", "err", err)
		}
	}()
	return nil
}

func (s *AccountServiceImpl) UpdateFee(ctx context.Context, accountID string) error {
	acc := domain.Account{ID: accountID}
	last, err := s.adjustmentRepo.GetLastByAccount(ctx, acc)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "applying last tariff adjustment", "adjustment", last)
	if err := s.accountRepo.UpdateFee(ctx, acc, last.NewFee); err != nil {
		return err
	}
	return nil
}

func (s *AccountServiceImpl) GetTariffAdjustments(ctx context.Context, acc domain.Account) ([]domain.TariffAdjustmentRequest, error) {
	return s.adjustmentRepo.AllByAccount(ctx, acc)
}

func (s *AccountServiceImpl) GetAccount(ctx context.Context, accountID string) (domain.Account, error) {
	return s.accountRepo.Get(ctx, domain.Account{ID: accountID})
}
