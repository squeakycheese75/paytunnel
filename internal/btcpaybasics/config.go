package btcpaybasics

import (
	"fmt"
	"os"
)

type Config struct {
	Port                string
	BTCPayWebhookSecret string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		Port:                envOrDefault("PORT", "8080"),
		BTCPayWebhookSecret: os.Getenv("BTCPAY_WEBHOOK_SECRET"),
	}

	if cfg.BTCPayWebhookSecret == "" {
		return Config{}, fmt.Errorf("BTCPAY_WEBHOOK_SECRET is required")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
