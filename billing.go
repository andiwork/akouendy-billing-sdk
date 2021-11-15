package billing

import (
	"crypto/sha512"
	"encoding/hex"
)

func ValidatePaymentWebhook(webhook PaymentWebhook) bool {
	strHash := billingConfig.AppToken + "|" + webhook.TransactionID + "|" + webhook.Status
	hashHex := hash512(strHash)
	if hashHex == webhook.Hash {
		return true
	}
	return false
}

func hash512(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
