package security

type Crypto interface {
	Decrypt(cipherText string) ([]byte, error)
	Encrypt(plaintext string) (string, error)
}
