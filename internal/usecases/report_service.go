package usecases

import (
	"cmp"
	"context"
	"slices"

	"sre/internal/domain"
)

var _ ReportService = (*reportService)(nil)

// ReportService builds fintech summary reports (totals by type, top by fee).
type ReportService interface {
	GetReport(ctx context.Context) (domain.Report, error)
}

type reportService struct {
	searchService SearchService
}

// NewReportService creates a ReportService.
func NewReportService(searchService SearchService) *reportService {
	return &reportService{searchService: searchService}
}

func (s *reportService) GetReport(ctx context.Context) (domain.Report, error) {
	accounts, err := s.searchService.SearchAccountsByTerm(ctx, "")
	if err != nil {
		return domain.Report{}, err
	}
	return domain.Report{
		TotalAccounts: s.totalAccounts(accounts),
		TotalsByType:  s.totalsByType(accounts),
		Top100ByFee:   s.top100ByFee(accounts),
	}, nil
}

func (s *reportService) totalAccounts(accounts []domain.Account) int {
	return len(accounts)
}

func (s *reportService) totalsByType(accounts []domain.Account) map[string]int {
	out := make(map[string]int)
	for _, a := range accounts {
		out[a.Type]++
	}
	return out
}

func (s *reportService) top100ByFee(accounts []domain.Account) []domain.Account {
	cp := make([]domain.Account, len(accounts))
	copy(cp, accounts)
	slices.SortFunc(cp, func(i, j domain.Account) int {
		return cmp.Compare(i.MonthlyFee, j.MonthlyFee)
	})
	slices.Reverse(cp)
	if len(cp) > 100 {
		return cp[:100]
	}
	return cp
}
