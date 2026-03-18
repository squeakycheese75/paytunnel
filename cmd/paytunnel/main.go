package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/squeakycheese75/paytunnel/internal/simulator"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "simulate":
		runSimulate(os.Args[2:])
	default:
		fmt.Println("unknown command:", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runSimulate(args []string) {
	fs := flag.NewFlagSet("simulate", flag.ExitOnError)

	url := fs.String("url", "http://localhost:8080/webhook/btcpay", "target webhook URL")
	secret := fs.String("secret", "my-supersecret-key", "BTCPay webhook secret")
	duplicate := fs.Int("duplicate", 1, "number of times to send the same event")
	invoiceID := fs.String("invoice-id", "inv_123", "invoice ID to include in the webhook payload")

	if err := fs.Parse(args); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	rest := fs.Args()
	if len(rest) < 1 {
		fmt.Println("usage: paytunnel simulate [--url ...] [--secret ...] [--duplicate N] <event>")
		os.Exit(1)
	}

	if *duplicate < 1 {
		fmt.Println("error: --duplicate must be at least 1")
		os.Exit(1)
	}

	event := rest[0]

	opts := simulator.Options{
		URL:       *url,
		Secret:    *secret,
		Duplicate: *duplicate,
		InvoiceID: *invoiceID,
	}

	if err := simulator.Simulate(event, opts); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  paytunnel simulate [--url ...] [--secret ...] [--duplicate N] <event>")
	fmt.Println()
	fmt.Println("example:")
	fmt.Println("  paytunnel simulate --url http://localhost:8080/webhook/btcpay --secret my-supersecret-key --duplicate 2 invoice.paid")
}
