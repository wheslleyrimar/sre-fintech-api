package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"sre/internal/domain"
	httpclient "sre/internal/httpClient"
	"sre/internal/usecases"
	"sre/internal/utils"
)

var _ usecases.AccountSearcher = (*SearchEngine)(nil)

// NewSearchEngine creates a client for the fintech search/accounts API.
func NewSearchEngine(factory httpclient.EndpointFactory) *SearchEngine {
	return &SearchEngine{
		getEndpoint: factory.Build("/v1/accounts"),
	}
}

type SearchEngine struct {
	getEndpoint httpclient.Endpoint
}

func (s *SearchEngine) SearchByTerm(ctx context.Context, term string) ([]domain.Account, error) {
	q := make(url.Values)
	if term != "" {
		q.Set("term", term)
	}
	var opts []httpclient.RequestOption
	if len(q) > 0 {
		opts = append(opts, httpclient.WithQuery(q))
	}
	res, err := s.getEndpoint.Get(ctx, opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, &httpclient.HTTPError{StatusCode: res.StatusCode, Body: fmt.Sprintf("search API returned %d: %s", res.StatusCode, string(body))}
	}
	var list []SearchResultAccount
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, err
	}
	out, err := utils.Map(list, func(r SearchResultAccount) domain.Account {
		return domain.Account{
			ID:         r.ID,
			Name:       r.Name,
			MonthlyFee: r.MonthlyFee,
			Type:       r.Type,
		}
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

type SearchResultAccount struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	MonthlyFee float64 `json:"monthly_fee"`
	Type       string  `json:"type"`
}
