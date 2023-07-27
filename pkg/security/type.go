package security

type Crypto interface {
	String() string
	Encrypt(plaintext string) (string, error)
	Decrypt(cipherText string) ([]byte, error)
}
