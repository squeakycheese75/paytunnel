package simulator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/squeakycheese75/paytunnel/internal/repository"
	"github.com/squeakycheese75/paytunnel/internal/signing"
)

type Options struct {
	URL       string
	Secret    string
	Duplicate int
	InvoiceID string
	Delay     time.Duration
}

type Repo interface {
	Create(ctx context.Context, event repository.Event) error
}

type Simulator struct {
	repo Repo
}

func NewSimulator(repo Repo) *Simulator {
	return &Simulator{repo: repo}
}

func (s *Simulator) Simulate(event string, opts Options) error {
	if opts.InvoiceID == "" {
		opts.InvoiceID = "inv_123"
	}

	switch event {
	case "invoice.paid":
		return s.sendInvoicePaid(opts)
	case "invoice.expired":
		return s.sendInvoiceExpired(opts)
	case "invoice.underpaid":
		return s.sendInvoiceUnderpaid(opts)
	default:
		return fmt.Errorf("unknown event: %s", event)
	}
}

func (s *Simulator) sendInvoicePaid(opts Options) error {
	payload := map[string]any{
		"deliveryId":   "d_" + randomID(),
		"webhookId":    "w_test",
		"isRedelivery": false,
		"type":         "InvoiceSettled",
		"timestamp":    time.Now().Unix(),
		"storeId":      "store_123",
		"data": map[string]any{
			"id":               opts.InvoiceID,
			"status":           "Settled",
			"additionalStatus": "paid",
		},
	}

	return s.sendPayload("invoice.paid", payload, opts)
}

func (s *Simulator) sendInvoiceExpired(opts Options) error {
	payload := map[string]any{
		"deliveryId":   "d_" + randomID(),
		"webhookId":    "w_test",
		"isRedelivery": false,
		"type":         "InvoiceExpired",
		"timestamp":    time.Now().Unix(),
		"storeId":      "store_123",
		"data": map[string]any{
			"id":               opts.InvoiceID,
			"status":           "Expired",
			"additionalStatus": "expired",
		},
	}

	return s.sendPayload("invoice.expired", payload, opts)
}

func (s *Simulator) sendInvoiceUnderpaid(opts Options) error {
	payload := map[string]any{
		"deliveryId":   "d_" + randomID(),
		"webhookId":    "w_test",
		"isRedelivery": false,
		"type":         "InvoicePaymentSettled",
		"timestamp":    time.Now().Unix(),
		"storeId":      "store_123",
		"data": map[string]any{
			"id":               opts.InvoiceID,
			"status":           "Settled",
			"additionalStatus": "underpaid",
		},
	}

	return s.sendPayload("invoice.underpaid", payload, opts)
}

func (s *Simulator) sendPayload(name string, payload map[string]any, opts Options) error {
	if opts.Duplicate < 1 {
		opts.Duplicate = 1
	}

	if opts.Delay > 0 {
		log.Printf("waiting %s before sending %s\n", opts.Delay, name)
		time.Sleep(opts.Delay)
	}

	deliveryID, _ := payload["deliveryId"].(string)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if err := s.repo.Create(context.Background(), repository.Event{
		DeliveryID: deliveryID,
		EventName:  name,
		TargetUrl:  opts.URL,
		BodyJson:   string(body),
		Secret:     opts.Secret,
		CreatedAt:  time.Now().UTC(),
	}); err != nil {
		return fmt.Errorf("save event: %w", err)
	}

	for i := 1; i <= opts.Duplicate; i++ {
		attemptPayload := clonePayload(payload)
		if i > 1 {
			attemptPayload["isRedelivery"] = true
		}

		attemptBody, err := json.Marshal(attemptPayload)
		if err != nil {
			return fmt.Errorf("marshal attempt payload: %w", err)
		}

		sig := signing.BTCPaySignature(attemptBody, opts.Secret)

		req, err := http.NewRequest(http.MethodPost, opts.URL, bytes.NewBuffer(attemptBody))
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("BTCPay-Sig", "sha256="+sig)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("send request: %w", err)
		}
		_ = resp.Body.Close()

		log.Printf(
			"\n[event] %s\ninvoice_id=%s\ndelivery_id=%s\nattempt=%d/%d\nstatus=%s\n\n",
			name,
			opts.InvoiceID,
			deliveryID,
			i,
			opts.Duplicate,
			resp.Status,
		)
	}

	return nil
}

func clonePayload(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func randomID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
