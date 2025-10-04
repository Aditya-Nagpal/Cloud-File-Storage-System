package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

var otpDigits = 6

func GenerateOtp() (string, error) {
	max := int64(1)
	for range otpDigits {
		max *= 10
	}

	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return "", err
	}
	format := fmt.Sprintf("%%0%dd", otpDigits)
	return fmt.Sprintf(format, n.Int64()), nil
}

func GenerateSalt(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HashOTP(otp, salt, pepper string) string {
	key := []byte(pepper)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(salt))
	h.Write([]byte(otp))
	return hex.EncodeToString(h.Sum(nil))
}
