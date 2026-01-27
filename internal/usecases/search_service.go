package usecases

import (
	"context"

	"sre/internal/domain"

	"log/slog"
)

var _ SearchService = (*SearchServiceImpl)(nil)

// SearchService searches accounts by term.
type SearchService interface {
	SearchAccountsByTerm(ctx context.Context, term string) ([]domain.Account, error)
}

// NewSearchService creates a SearchService.
func NewSearchService(searcher AccountSearcher) *SearchServiceImpl {
	return &SearchServiceImpl{searcher: searcher}
}

type SearchServiceImpl struct {
	searcher AccountSearcher
}

func (s *SearchServiceImpl) SearchAccountsByTerm(ctx context.Context, term string) ([]domain.Account, error) {
	accounts, err := s.searcher.SearchByTerm(ctx, term)
	if err != nil {
		slog.ErrorContext(ctx, "search by term failed", "term", term, "err", err)
		return nil, err
	}
	return accounts, nil
}
