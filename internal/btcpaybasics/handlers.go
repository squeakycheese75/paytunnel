package btcpaybasics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (a *App) healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (a *App) ordersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(a.orderStore.List()); err != nil {
		http.Error(w, "failed to encode orders", http.StatusInternalServerError)
		return
	}
}

func (a *App) btcpayWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	if err := VerifyBTCPaySignature(r.Header.Get("BTCPay-Sig"), rawBody, a.config.BTCPayWebhookSecret); err != nil {
		http.Error(w, fmt.Sprintf("invalid signature: %v", err), http.StatusUnauthorized)
		return
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, rawBody, "", "  "); err == nil {
		log.Printf("raw webhook body:\n%s", pretty.String())
	}

	var event BTCPayWebhook
	if err := json.Unmarshal(rawBody, &event); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}

	if event.DeliveryID == "" {
		http.Error(w, "missing deliveryId", http.StatusBadRequest)
		return
	}

	if !a.deliveryStore.MarkIfNew(event.DeliveryID) {
		log.Printf(
			"duplicate btcpay webhook ignored: type=%q delivery_id=%q webhook_id=%q redelivery=%t",
			event.Type,
			event.DeliveryID,
			event.WebhookID,
			event.IsRedelivery,
		)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("duplicate ignored"))
		return
	}

	var invoice BTCPayInvoiceData
	if len(event.Data) > 0 {
		_ = json.Unmarshal(event.Data, &invoice)
	}

	log.Printf(
		"received btcpay webhook: type=%q delivery_id=%q webhook_id=%q redelivery=%t invoice_id=%q status=%q additional_status=%q",
		event.Type,
		event.DeliveryID,
		event.WebhookID,
		event.IsRedelivery,
		invoice.ID,
		invoice.Status,
		invoice.AdditionalStatus,
	)

	if event.Type == "InvoiceSettled" && invoice.ID != "" {
		order, found := a.orderStore.MarkPaidByInvoiceID(invoice.ID)
		if !found {
			log.Printf("no order found for invoice_id=%q", invoice.ID)
		} else {
			log.Printf("order %q marked as paid for invoice_id=%q", order.ID, invoice.ID)
		}
	}

	if event.Type == "InvoiceExpired" && invoice.ID != "" {
		log.Printf("invoice %q expired", invoice.ID)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
