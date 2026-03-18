package simulator

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Options struct {
	URL       string
	Secret    string
	Duplicate int
	InvoiceID string
}

func Simulate(event string, opts Options) error {
	switch event {
	case "invoice.paid":
		return sendInvoicePaid(opts)
	case "invoice.expired":
		return sendInvoiceExpired(opts)
	default:
		return fmt.Errorf("unknown event: %s", event)
	}
}

func sendInvoicePaid(opts Options) error {
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

	return sendPayload("invoice.paid", payload, opts)
}

func sendInvoiceExpired(opts Options) error {
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

	return sendPayload("invoice.expired", payload, opts)
}

func sendPayload(name string, payload map[string]any, opts Options) error {
	if opts.Duplicate < 1 {
		opts.Duplicate = 1
	}

	if opts.InvoiceID == "" {
		opts.InvoiceID = "inv_123"
	}

	deliveryID, _ := payload["deliveryId"].(string)

	for i := 1; i <= opts.Duplicate; i++ {
		if i > 1 {
			payload["isRedelivery"] = true
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal payload: %w", err)
		}

		sig := sign(body, opts.Secret)

		req, err := http.NewRequest(http.MethodPost, opts.URL, bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("BTCPay-Sig", "sha256="+sig)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("send request: %w", err)
		}
		resp.Body.Close()

		fmt.Printf(
			"send %d/%d %s delivery_id=%s -> status: %s\n",
			i,
			opts.Duplicate,
			name,
			deliveryID,
			resp.Status,
		)
	}

	return nil
}

func sign(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func randomID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
