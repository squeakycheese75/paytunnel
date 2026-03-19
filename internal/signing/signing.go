package signing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func BTCPaySignature(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
