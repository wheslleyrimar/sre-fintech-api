package integrations

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"sre/internal/domain"
	httpclient "sre/internal/httpClient"
	"sre/internal/usecases"
	"sre/internal/utils"
)

var (
	_ usecases.AccountRepository        = (*AccountsApi)(nil)
	_ usecases.TariffAdjustmentRepository = (*AccountsApi)(nil)
)

// NewAccountsApi creates an HTTP client for the fintech accounts API.
func NewAccountsApi(factory httpclient.EndpointFactory) *AccountsApi {
	return &AccountsApi{
		accountEndpoint:            factory.Build("/v1/accounts/{id}"),
		adjustmentsEndpoint:       factory.Build("/v1/accounts/{id}/tariff-adjustments"),
		adjustmentsLastEndpoint:   factory.Build("/v1/accounts/{id}/tariff-adjustments/last"),
		postAdjustmentEndpoint:   factory.Build("/v1/accounts/{id}/tariff-adjustments"),
	}
}

type AccountsApi struct {
	accountEndpoint          httpclient.Endpoint
	adjustmentsEndpoint      httpclient.Endpoint
	adjustmentsLastEndpoint  httpclient.Endpoint
	postAdjustmentEndpoint  httpclient.Endpoint
}

func (a *AccountsApi) GetLastByAccount(ctx context.Context, acc domain.Account) (*domain.TariffAdjustmentRequest, error) {
	res, err := a.adjustmentsLastEndpoint.Get(ctx, httpclient.WithParam("id", acc.ID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, &httpclient.HTTPError{StatusCode: res.StatusCode, Body: string(body)}
	}
	var r LastAdjustmentResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &domain.TariffAdjustmentRequest{
		TransactionID: r.TransactionID,
		AccountID:     r.AccountID,
		NewFee:        r.Fee,
	}, nil
}

func (a *AccountsApi) AllByAccount(ctx context.Context, acc domain.Account) ([]domain.TariffAdjustmentRequest, error) {
	res, err := a.adjustmentsEndpoint.Get(ctx, httpclient.WithParam("id", acc.ID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, &httpclient.HTTPError{StatusCode: res.StatusCode, Body: string(body)}
	}
	var list []AdjustmentResponse
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, err
	}
	out, err := utils.Map(list, func(r AdjustmentResponse) domain.TariffAdjustmentRequest {
		return domain.TariffAdjustmentRequest{
			TransactionID: r.TransactionID,
			AccountID:     r.AccountID,
			NewFee:        r.Fee,
		}
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (a *AccountsApi) Create(ctx context.Context, input domain.TariffAdjustmentRequest, callbackURL string) error {
	body := CreateAdjustmentBody{NewFee: input.NewFee, CallbackURL: callbackURL}
	res, err := a.postAdjustmentEndpoint.Post(ctx,
		httpclient.WithParam("id", input.AccountID),
		httpclient.WithBody(body),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return &httpclient.HTTPError{StatusCode: res.StatusCode, Body: string(b)}
	}
	return nil
}

func (a *AccountsApi) UpdateFee(ctx context.Context, acc domain.Account, newFee float64) error {
	res, err := a.accountEndpoint.Patch(ctx,
		httpclient.WithParam("id", acc.ID),
		httpclient.WithBody(UpdateAccountBody{MonthlyFee: newFee}),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return &httpclient.HTTPError{StatusCode: res.StatusCode, Body: string(b)}
	}
	return nil
}

func (a *AccountsApi) Get(ctx context.Context, id domain.Account) (domain.Account, error) {
	res, err := a.accountEndpoint.Get(ctx, httpclient.WithParam("id", id.ID))
	if err != nil {
		return domain.Account{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return domain.Account{}, &httpclient.HTTPError{StatusCode: res.StatusCode, Body: string(b)}
	}
	var r AccountResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return domain.Account{}, err
	}
	return domain.Account{
		ID:         r.ID,
		Name:       r.Name,
		MonthlyFee: r.MonthlyFee,
		Type:       r.Type,
	}, nil
}

type AccountResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	MonthlyFee float64 `json:"monthly_fee"`
	Type       string  `json:"type"`
}

type AdjustmentResponse struct {
	TransactionID string  `json:"transaction_id"`
	Fee           float64 `json:"fee"`
	AccountID     string  `json:"account_id"`
}

type LastAdjustmentResponse struct {
	TransactionID string  `json:"transaction_id"`
	Fee           float64 `json:"fee"`
	AccountID     string  `json:"account_id"`
}

type CreateAdjustmentBody struct {
	NewFee      float64 `json:"new_fee"`
	CallbackURL string  `json:"callback_url"`
}

type UpdateAccountBody struct {
	MonthlyFee float64 `json:"monthly_fee"`
}
