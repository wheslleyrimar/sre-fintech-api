package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"sre/internal/domain"
	"sre/internal/usecases"
	"sre/internal/utils"
)

// NewAccountController creates an account controller.
func NewAccountController(s usecases.AccountService) *AccountController {
	return &AccountController{usecase: s}
}

type AccountController struct {
	usecase usecases.AccountService
}

// Routes registers account and tariff-adjustment routes on r.
func (c *AccountController) Routes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", c.listAccounts)
		r.Get("/{id}", c.getAccount)
		r.Post("/{id}/tariff-adjustments", c.createTariffAdjustment)
		r.Get("/{id}/tariff-adjustments", c.getTariffAdjustments)
		r.Post("/notifications", c.notifications)
	})
}

func (c *AccountController) listAccounts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (c *AccountController) getAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		encodeError(w, "missing id", http.StatusBadRequest)
		return
	}
	acc, err := c.usecase.GetAccount(r.Context(), id)
	if err != nil {
		encodeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encodeJSON(w, GetAccountResponse{
		ID:         acc.ID,
		Name:       acc.Name,
		MonthlyFee: acc.MonthlyFee,
		Type:       acc.Type,
	}, http.StatusOK)
}

func (c *AccountController) createTariffAdjustment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		encodeError(w, "missing id", http.StatusBadRequest)
		return
	}
	var payload TariffAdjustmentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		encodeError(w, "invalid body", http.StatusBadRequest)
		return
	}
	input := domain.TariffAdjustmentRequest{
		AccountID:     id,
		TransactionID: uuid.NewString(),
		NewFee:        payload.NewFee,
	}
	if err := c.usecase.SendTariffAdjustmentRequest(r.Context(), input); err != nil {
		encodeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *AccountController) getTariffAdjustments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		encodeError(w, "missing id", http.StatusBadRequest)
		return
	}
	acc := domain.Account{ID: id}
	list, err := c.usecase.GetTariffAdjustments(r.Context(), acc)
	if err != nil {
		encodeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := utils.Map(list, func(a domain.TariffAdjustmentRequest) TariffAdjustmentResponse {
		return TariffAdjustmentResponse{
			TransactionID: a.TransactionID,
			AccountID:     a.AccountID,
			NewFee:        a.NewFee,
		}
	})
	if err != nil {
		encodeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encodeJSON(w, out, http.StatusOK)
}

func (c *AccountController) notifications(w http.ResponseWriter, r *http.Request) {
	var msg NotificationMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		encodeError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := c.usecase.UpdateFee(r.Context(), msg.AccountID); err != nil {
		encodeError(w, "update fee failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

type GetAccountResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	MonthlyFee float64 `json:"monthly_fee"`
	Type       string  `json:"type"`
}

type TariffAdjustmentPayload struct {
	NewFee float64 `json:"new_fee"`
}

type TariffAdjustmentResponse struct {
	TransactionID string  `json:"transaction_id"`
	AccountID     string  `json:"account_id"`
	NewFee        float64 `json:"new_fee"`
}

type NotificationMessage struct {
	TransactionID string `json:"transaction_id"`
	AccountID     string `json:"account_id"`
	Status        string `json:"status"`
}
