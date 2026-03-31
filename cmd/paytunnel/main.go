package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/squeakycheese75/paytunnel/internal/db"
	"github.com/squeakycheese75/paytunnel/internal/eventlog"
	"github.com/squeakycheese75/paytunnel/internal/relay"
	"github.com/squeakycheese75/paytunnel/internal/repository"
	"github.com/squeakycheese75/paytunnel/internal/simulator"
	"github.com/squeakycheese75/paytunnel/internal/tunnel"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	sqlDB, err := db.NewDB()
	if err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Println("close db error:", err)
		}
	}()

	repo := repository.NewEventRepository(sqlDB)

	switch os.Args[1] {
	case "simulate":
		runSimulate(os.Args[2:], repo)
	case "events":
		runEvents(os.Args[2:], repo)
	case "relay":
		runRelay(os.Args[2:])
	case "tunnel":
		runTunnel(os.Args[2:])
	default:
		log.Println("unknown command:", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runSimulate(args []string, repo *repository.EventRepository) {
	fs := flag.NewFlagSet("simulate", flag.ExitOnError)

	url := fs.String("url", "http://localhost:8080/webhook/btcpay", "target webhook URL")
	secret := fs.String("secret", "my-supersecret-key", "BTCPay webhook secret")
	duplicate := fs.Int("duplicate", 1, "number of times to send the same event")
	invoiceID := fs.String("invoice-id", "inv_123", "invoice ID to include in the webhook payload")
	delay := fs.Duration("delay", 0, "delay before sending the webhook, e.g. 2s")

	if err := fs.Parse(args); err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}

	rest := fs.Args()
	if len(rest) < 1 {
		log.Println("usage: paytunnel simulate [--url ...] [--secret ...] [--duplicate N] <event>")
		os.Exit(1)
	}

	if *duplicate < 1 {
		log.Println("error: --duplicate must be at least 1")
		os.Exit(1)
	}

	event := rest[0]

	opts := simulator.Options{
		URL:       *url,
		Secret:    *secret,
		Duplicate: *duplicate,
		InvoiceID: *invoiceID,
		Delay:     *delay,
	}

	sim := simulator.NewSimulator(repo)

	if err := sim.Simulate(event, opts); err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
}

func runEvents(args []string, repo *repository.EventRepository) {
	if len(args) < 1 {
		log.Println("usage: paytunnel events <list|replay>")
		os.Exit(1)
	}

	eventlog := eventlog.NewEventLog(repo)

	switch args[0] {
	case "list":
		events, err := eventlog.List(context.Background())
		if err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}

		for _, e := range events {
			log.Printf("%s  %s  %s\n", e.DeliveryID, e.EventName, e.CreatedAt)
		}

	case "replay":
		if len(args) < 2 {
			log.Println("usage: paytunnel events replay <delivery-id>")
			os.Exit(1)
		}

		if err := eventlog.ReplayEvent(context.Background(), args[1]); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}

	default:
		log.Println("unknown events command:", args[0])
		os.Exit(1)
	}
}

func runRelay(args []string) {
	fs := flag.NewFlagSet("relay", flag.ExitOnError)
	listen := fs.String("listen", ":9000", "address for relay to listen on")

	if err := fs.Parse(args); err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}

	server := relay.NewServer(*listen)
	if err := server.Run(); err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
}

func runTunnel(args []string) {
	fs := flag.NewFlagSet("tunnel", flag.ExitOnError)
	relayURL := fs.String("relay", "ws://localhost:9000/ws", "relay websocket URL")
	targetURL := fs.String("to", "http://localhost:8080", "local target base URL")

	if err := fs.Parse(args); err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}

	client := tunnel.NewClient(*relayURL, *targetURL)
	if err := client.Run(); err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
}

func printUsage() {
	log.Println("usage:")
	log.Println("  paytunnel simulate [--url ...] [--secret ...] [--duplicate N] <event>")
	log.Println()
	log.Println("example:")
	log.Println("  paytunnel simulate --url http://localhost:8080/webhook/btcpay --secret my-supersecret-key --duplicate 2 invoice.paid")
}
