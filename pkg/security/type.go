package security

import "net/http"

type Crypto interface {
	String() string
	Encrypt(plaintext string) (string, error)
	Decrypt(cipherText string) ([]byte, error)
}

type Verify interface {
	VerifyRequest(r *http.Request) error
}

// SignatureStrategy 抽象签名算法策略接口（支持 HMAC、RSA 等）
type SignatureStrategy interface {
	Generate(signingString, secret string) string
	Verify(signingString, secret, givenSig string) bool
}
