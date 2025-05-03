package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type HMACStrategy struct{}

func (s HMACStrategy) Generate(signingString, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signingString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s HMACStrategy) Verify(signingString, secret, givenSig string) bool {
	expected := s.Generate(signingString, secret)
	return hmac.Equal([]byte(expected), []byte(givenSig))
}
