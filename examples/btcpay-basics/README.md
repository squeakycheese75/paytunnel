# btcpay-basics

Minimal Go example showing how to **receive, verify, and safely handle BTCPay webhooks**.

This example demonstrates:

* ✅ Receiving a webhook
* 🔐 Verifying the `BTCPay-Sig` HMAC signature
* 🔁 Handling duplicate deliveries (idempotency)
* 🧾 Parsing invoice data
* 🪵 Logging useful event details

---

## 🚀 Run the example

From the repo root:

```bash
go run examples/btcpay-basics/main.go
```

The server will start on:

```text
http://localhost:8080
```

Endpoints:

* `GET /health`
* `POST /webhook/btcpay`

---

## ⚙️ Configuration

Create a `.env` file (or use environment variables):

```bash
PORT=8080
BTCPAY_WEBHOOK_SECRET=my-supersecret-key
```

---

## 🧪 Send a test webhook (signed)

### 1. Define payload

```bash
payload='{"deliveryId":"d_123","webhookId":"w_123","originalDeliveryId":null,"isRedelivery":false,"type":"InvoiceSettled","timestamp":1710000000,"storeId":"store_123","data":{"id":"inv_123","status":"Settled","additionalStatus":"paid"}}'
```

---

### 2. Generate signature

```bash
sig=$(printf '%s' "$payload" | openssl dgst -sha256 -hmac 'my-supersecret-key' -hex | sed 's/^.* //')
```

---

### 3. Send request

```bash
curl -i \
  -X POST http://localhost:8080/webhook/btcpay \
  -H "Content-Type: application/json" \
  -H "BTCPay-Sig: sha256=$sig" \
  -d "$payload"
```

---

## ✅ Expected output

You should see:

```text
received btcpay webhook: type="InvoiceSettled" delivery_id="d_123" ...
```

---

## 🔁 Duplicate delivery test

Send the **same request again**.

You should see:

```text
duplicate btcpay webhook ignored: delivery_id="d_123"
```

This demonstrates **idempotent webhook handling**, which is critical in production systems.

---

## ❌ Invalid signature test

```bash
curl -i \
  -X POST http://localhost:8080/webhook/btcpay \
  -H "Content-Type: application/json" \
  -H "BTCPay-Sig: sha256=invalid" \
  -d "$payload"
```

Expected:

```text
401 Unauthorized
```

---

## 🧠 Key concepts

### Webhook signature verification

BTCPay signs each webhook using:

```text
HMAC-SHA256(payload, webhook_secret)
```

The result is sent in:

```text
BTCPay-Sig: sha256=<signature>
```

---

### Idempotency

Webhook providers may:

* retry failed deliveries
* send duplicates
* deliver events out of order

This example:

* tracks processed `deliveryId`
* ignores duplicates safely
* always returns `200 OK` for already-processed events

---

## 🔗 Using with BTCPay Server

In BTCPay:

1. Go to:

   ```
   Stores → Settings → Webhooks
   ```
2. Create a webhook:

```text
URL: http://localhost:8080/webhook/btcpay
Secret: my-supersecret-key
Events: InvoiceSettled
```

---

## 🧭 Next steps

* Add persistent idempotency (database or cache)
* Handle order/payment state transitions
* Simulate webhook events locally
* Add retry handling and backoff testing

---

## 💡 Why this example matters

Most webhook examples only show how to receive requests.

This example shows how to do it **correctly in production**:

* verify authenticity
* prevent duplicate processing
* extract useful business data

---

## 📁 Related examples (coming next)

* `btcpay-order-status` — update order state from webhook
* `btcpay-retries` — simulate failures and retries
* `btcpay-simulator` — generate test events locally
