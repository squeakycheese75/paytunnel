package btcpaybasics

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
)

func VerifyBTCPaySignature(header string, body []byte, secret string) error {
	if header == "" {
		return fmt.Errorf("missing BTCPay-Sig header")
	}

	const prefix = "sha256="
	if !strings.HasPrefix(header, prefix) {
		return fmt.Errorf("signature header missing %q prefix", prefix)
	}

	gotHex := strings.TrimPrefix(header, prefix)
	gotSig, err := hex.DecodeString(gotHex)
	if err != nil {
		return fmt.Errorf("signature is not valid hex: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	expectedSig := mac.Sum(nil)

	if subtle.ConstantTimeCompare(gotSig, expectedSig) != 1 {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}
