package security

type Crypto interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(cipherText string) ([]byte, error)
}
