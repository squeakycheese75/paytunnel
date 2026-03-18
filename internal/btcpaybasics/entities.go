package btcpaybasics

import "encoding/json"

type BTCPayWebhook struct {
	DeliveryID         string          `json:"deliveryId"`
	WebhookID          string          `json:"webhookId"`
	OriginalDeliveryID *string         `json:"originalDeliveryId"`
	IsRedelivery       bool            `json:"isRedelivery"`
	Type               string          `json:"type"`
	Timestamp          int64           `json:"timestamp"`
	StoreID            string          `json:"storeId"`
	Data               json.RawMessage `json:"data"`
}

type BTCPayInvoiceData struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	AdditionalStatus string `json:"additionalStatus"`
}

type Order struct {
	ID        string `json:"id"`
	InvoiceID string `json:"invoiceId"`
	Status    string `json:"status"`
}
