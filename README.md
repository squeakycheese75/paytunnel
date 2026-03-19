# paytunnel

A lightweight CLI for simulating and replaying payment webhooks (starting with BTCPay).

## Why?

Testing payment flows locally is painful:

* Hard to trigger real payment scenarios
* Difficult to replay webhooks
* No easy way to simulate retries or edge cases

`paytunnel` helps you:

* Simulate payment webhooks locally
* Store webhook events
* Replay events on demand

---

## Quick start

### 1. Build

```bash
go build -o bin/paytunnel ./cmd/paytunnel
```

---

### 2. Run example server

```bash
go run ./examples/btcpay-basics
```

This starts a local webhook receiver:

```
http://localhost:8080/webhook/btcpay
```

---

### 3. Simulate a payment

```bash
./bin/paytunnel simulate \
  --url http://localhost:8080/webhook/btcpay \
  --secret my-supersecret-key \
  invoice.paid
```

---

### 4. List events

```bash
./bin/paytunnel events list
```

---

### 5. Replay an event

```bash
./bin/paytunnel events replay <delivery-id>
```

---

## Supported scenarios

* invoice.paid
* invoice.expired
* invoice.underpaid

---

## Features

* 🔁 Replay webhooks
* 🧪 Simulate payment scenarios
* 💾 Local event storage (SQLite)
* 🔐 BTCPay-compatible signature verification

---

## Project structure

```
internal/
  repository/   # DB access
  simulator/    # event generation + sending
  eventlog/     # list + replay
```

---

## Roadmap

* More payment scenarios (overpaid, late, retries)
* Webhook retry simulation
* Out-of-order delivery
* HTTP tunneling (ngrok-style)

---

## Development

Run example:

```bash
go run ./examples/btcpay-basics
```

Simulate:

```bash
./bin/paytunnel simulate --url http://localhost:8080/webhook/btcpay --secret my-supersecret-key invoice.paid
```

---

## License

MIT
