# paytunnel

Simulate Bitcoin / BTCPay webhook events locally.

Test your payment logic without running BTCPay Server or using real Bitcoin.

---

## 🚀 Quick start

### 1. Run the example server

```bash
go run ./examples/btcpay-basics
```

### 2. Simulate a payment

```bash
go run ./cmd/paytunnel simulate invoice.paid
```

### 3. Check orders

```bash
curl http://localhost:8080/orders
```

---

## ✨ Features

* Simulate BTCPay webhook events
* Signature verification (HMAC SHA256)
* Duplicate delivery handling
* Multiple scenarios:

  * `invoice.paid`
  * `invoice.expired`
* Configurable:

  * `--url`
  * `--secret`
  * `--duplicate`
  * `--invoice-id`

---

## 📦 Example

```bash
go run ./cmd/paytunnel simulate \
  --url http://localhost:8080/webhook/btcpay \
  --secret my-supersecret-key \
  --duplicate 2 \
  --invoice-id inv_123 \
  invoice.paid
```

---

## 🧪 What this does

1. Sends a signed BTCPay-style webhook
2. Your server verifies the signature
3. Processes the event
4. Updates application state (orders)

---

## 🧠 Why?

Testing payment flows is painful:

* Requires BTCPay setup
* Hard to simulate edge cases
* Manual curl + signature generation

**paytunnel** solves this with a simple CLI.

---

## 📁 Project structure

```text
cmd/paytunnel         # CLI entrypoint
internal/simulator    # webhook simulator
internal/btcpay       # shared models + signature logic
examples/btcpay-basics # runnable example server
```

---

## 🔮 Roadmap

* [ ] Delay / retry simulation (`--delay`)
* [ ] More scenarios (underpaid, failed)
* [ ] CLI install (`go install`)
* [ ] Webhook inspector
* [ ] Optional tunneling support

---

## 🤝 Contributing

PRs welcome. Ideas and feedback encouraged.

---

## ⭐️ Support

If you find this useful, consider giving it a star.
