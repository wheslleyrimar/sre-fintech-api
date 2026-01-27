package integrations

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"sre/internal/domain"
	httpclient "sre/internal/httpClient"
	"sre/internal/usecases"
)

var _ usecases.AdjustmentFlowProcessor = (*AdjustmentFlowProcessor)(nil)

// NewAdjustmentFlowProcessor creates a client for the fintech adjustment-approval-flow API.
func NewAdjustmentFlowProcessor(factory httpclient.EndpointFactory) *AdjustmentFlowProcessor {
	return &AdjustmentFlowProcessor{
		endpoint: factory.Build("/v1/adjustment-approval-flow"),
	}
}

type AdjustmentFlowProcessor struct {
	endpoint httpclient.Endpoint
}

func (p *AdjustmentFlowProcessor) BeginFlow(ctx context.Context, input domain.TariffAdjustmentRequest, callbackURL string) error {
	body := AdjustmentApprovalFlowBody{
		AccountID:   input.AccountID,
		NewFee:      input.NewFee,
		CallbackURL: callbackURL,
	}
	res, err := p.endpoint.Post(ctx, httpclient.WithBody(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("adjustment-approval-flow returned %d: %s", res.StatusCode, string(b))
	}
	return nil
}

type AdjustmentApprovalFlowBody struct {
	AccountID   string  `json:"account_id"`
	NewFee      float64 `json:"new_fee"`
	CallbackURL string  `json:"callback_url"`
}
