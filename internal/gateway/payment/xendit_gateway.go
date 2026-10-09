package payment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type XenditGateway struct {
	secretKey    string
	webhookToken string
	httpClient   *http.Client
}

type XenditInvoiceRequest struct {
	ExternalID      string `json:"external_id"`
	Amount          int64  `json:"amount"`
	PayerEmail      string `json:"payer_email"`
	Description     string `json:"description"`
	InvoiceDuration int    `json:"invoice_duration"` // Seconds
}

type XenditInvoiceResponse struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	InvoiceURL string `json:"invoice_url"`
	Status     string `json:"status"`
	Amount     int64  `json:"amount"`
}

func NewXenditGateway(viper *viper.Viper) *XenditGateway {
	return &XenditGateway{
		secretKey:    viper.GetString("XENDIT_SECRET_KEY"),
		webhookToken: viper.GetString("XENDIT_WEBHOOK_TOKEN"),
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (g *XenditGateway) CreateInvoice(ctx context.Context, req XenditInvoiceRequest) (*XenditInvoiceResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.xendit.co/v2/invoices", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	rawAuth := fmt.Appendf(nil, "%s:", g.secretKey)
	authHeader := base64.StdEncoding.EncodeToString(rawAuth)
	httpReq.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("xendit invoice creation failed with status: %d", res.StatusCode)
	}

	var response XenditInvoiceResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (g *XenditGateway) VerifyWebhookToken(token string) bool {
	if g.webhookToken == "" {
		return true // Allow development inspection if token check is unconfigured
	}
	return g.webhookToken == token
}
