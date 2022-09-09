package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/url"
)

const (
	CryptoRSAName = "rsa"
)

type RsaCiphers struct {
	ClientCipher     *RsaCipher
	ServerCipher     *RsaCipher
	ClientPublicKey  string
	ClientPrivateKey string
	ServerPublicKey  string
	ServerPrivateKey string
}

func NewRsaCiphers(serverPublicKey, serverPrivateKey, clientPublicKey, clientPrivateKey string) (*RsaCiphers, error) {
	var err error
	c := &RsaCiphers{
		ServerPublicKey:  serverPublicKey,
		ServerPrivateKey: serverPrivateKey,
		ClientPublicKey:  clientPublicKey,
		ClientPrivateKey: clientPrivateKey,
	}
	c.ClientCipher, err = NewRsaCipher(c.ClientPublicKey, c.ClientPrivateKey)
	if err != nil {
		return nil, err
	}
	c.ServerCipher, err = NewRsaCipher(c.ServerPublicKey, c.ServerPrivateKey)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (a *RsaCiphers) String() string {
	return CryptoRSAName
}

func (a *RsaCiphers) Encrypt(plainText string) (cryptText string, err error) {
	return a.ServerCipher.Encrypt(plainText)
}

func (a *RsaCiphers) Decrypt(cryptText string) (plainText []byte, err error) {
	return a.ClientCipher.Decrypt(cryptText)
}

type RsaCipher struct {
	PubKey     string
	PriKey     string
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewRsaCipher(pubKey, priKey string) (*RsaCipher, error) {
	c := &RsaCipher{
		PubKey: pubKey,
		PriKey: priKey,
	}
	// PublicKey
	block, _ := pem.Decode([]byte(c.PubKey))
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	c.PublicKey = publicKeyInterface.(*rsa.PublicKey)
	// PrivateKey
	block, _ = pem.Decode([]byte(c.PriKey))
	c.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (a *RsaCipher) String() string {
	return CryptoRSAName
}

func split(buf []byte, lim int) [][]byte {
	//glog.Debug(context.Background(), "长度:", len(buf))
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}

func (a *RsaCipher) Encrypt(plainText string) (cryptText string, err error) {
	//maxEncryptBlockLen := publicKey.N.BitLen()/8 - 11
	maxEncryptBlockLen := ((a.PublicKey.N.BitLen() + 7) >> 3) - 11
	chunks := split([]byte(a.EncodeURIComponent(plainText)), maxEncryptBlockLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		byteCode, err := rsa.EncryptPKCS1v15(rand.Reader, a.PublicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(byteCode)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func (a *RsaCipher) Decrypt(cryptText string) (plainText []byte, err error) {
	//maxDecryptBlockLen := privateKey.N.BitLen() / 8
	maxDecryptBlockLen := (a.PrivateKey.N.BitLen() + 7) >> 3
	decodeString, err := base64.StdEncoding.DecodeString(cryptText)
	if err != nil {
		return []byte{}, err
	}
	chunks := split(decodeString, maxDecryptBlockLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, a.PrivateKey, chunk)
		if err != nil {
			return []byte{}, err
		}
		buffer.Write(decrypted)
	}
	str := a.DecodeURIComponent(string(buffer.Bytes()))
	return []byte(str), err
}

func (a *RsaCipher) EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	//r = strings.Replace(r, "+", "%20", -1)
	return r
}

func (a *RsaCipher) DecodeURIComponent(str string) string {
	r, err := url.QueryUnescape(str)
	if err != nil {
		return str
	}
	return r
}
