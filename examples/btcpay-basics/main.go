package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/squeakycheese75/paytunnel/internal/btcpaybasics"
)

func main() {
	_ = godotenv.Load("examples/btcpay-basics/.env")

	cfg, err := btcpaybasics.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := btcpaybasics.NewApp(cfg)

	log.Fatal(app.Run())
}
