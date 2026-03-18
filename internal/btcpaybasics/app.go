package btcpaybasics

import (
	"fmt"
	"log"
	"net/http"
)

type App struct {
	config        Config
	deliveryStore *DeliveryStore
	orderStore    *OrderStore
}

func NewApp(cfg Config) *App {
	return &App{
		config:        cfg,
		deliveryStore: NewDeliveryStore(),
		orderStore:    NewOrderStore(),
	}
}

func (a *App) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.healthHandler)
	mux.HandleFunc("/orders", a.ordersHandler)
	mux.HandleFunc("/webhook/btcpay", a.btcpayWebhookHandler)

	addr := ":" + a.config.Port

	log.Printf("btcpay-basics listening on http://localhost%s", addr)
	log.Printf("health endpoint:  GET  http://localhost%s/health", addr)
	log.Printf("orders endpoint:  GET  http://localhost%s/orders", addr)
	log.Printf("webhook endpoint: POST http://localhost%s/webhook/btcpay", addr)

	return http.ListenAndServe(addr, mux)
}

func (a *App) address() string {
	return fmt.Sprintf(":%s", a.config.Port)
}
