package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"sre/internal/domain"
	"sre/internal/usecases"
	"sre/internal/utils"
)

// NewSearchController creates a search controller.
func NewSearchController(s usecases.SearchService) *SearchController {
	return &SearchController{usecase: s}
}

type SearchController struct {
	usecase usecases.SearchService
}

// Routes registers search routes on r.
func (c *SearchController) Routes(r chi.Router) {
	r.Get("/search", c.search)
}

func (c *SearchController) search(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	accounts, err := c.usecase.SearchAccountsByTerm(r.Context(), term)
	if err != nil {
		encodeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := utils.Map(accounts, func(a domain.Account) SearchResultItem {
		return SearchResultItem{
			ID:         a.ID,
			Name:       a.Name,
			MonthlyFee: a.MonthlyFee,
			Type:       a.Type,
		}
	})
	if err != nil {
		encodeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encodeJSON(w, SearchResponse{Data: out}, http.StatusOK)
}

type SearchResultItem struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	MonthlyFee float64 `json:"monthly_fee"`
	Type       string  `json:"type"`
}

type SearchResponse struct {
	Data []SearchResultItem `json:"data"`
}
