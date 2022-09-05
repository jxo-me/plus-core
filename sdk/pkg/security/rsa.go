package security

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/gogf/gf/v2/os/glog"
	"net/url"
	"runtime"
	"strings"
)

type RsaCipher struct {
	PubKey string
	PriKey string
}

func NewRsaCipher(pubKey, priKey string) (*RsaCipher, error) {
	c := &RsaCipher{
		PubKey: pubKey,
		PriKey: priKey,
	}
	return c, nil
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
	block, _ := pem.Decode([]byte(a.PubKey))
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				glog.Error(context.Background(), "runtime err:", err, "Check that the key is correct")
			default:
				glog.Error(context.Background(), "Rsa Encrypt error:", err)
			}
		}
	}()
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//maxEncryptBlockLen := publicKey.N.BitLen()/8 - 11
	maxEncryptBlockLen := ((publicKey.N.BitLen() + 7) >> 3) - 11
	chunks := split([]byte(a.EncodeURIComponent(plainText)), maxEncryptBlockLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		byteCode, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(byteCode)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func (a *RsaCipher) Decrypt(cryptText string) (plainText []byte, err error) {
	block, _ := pem.Decode([]byte(a.PriKey))
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				glog.Error(context.Background(), "runtime err:", err, "Check that the key is correct")
			default:
				glog.Error(context.Background(), "Rsa Decrypt error:", err)
			}
		}
	}()

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		glog.Error(context.Background(), "ParsePKCS1PrivateKey error:", err)
		return []byte{}, err
	}
	//maxDecryptBlockLen := privateKey.N.BitLen() / 8
	maxDecryptBlockLen := (privateKey.N.BitLen() + 7) >> 3
	decodeString, err := base64.StdEncoding.DecodeString(cryptText)
	if err != nil {
		glog.Error(context.Background(), "base64 decodeString error:", err)
		return []byte{}, err
	}
	chunks := split(decodeString, maxDecryptBlockLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk)
		if err != nil {
			glog.Error(context.Background(), "DecryptPKCS1v15 error:", err)
			return []byte{}, err
		}
		buffer.Write(decrypted)
	}
	return buffer.Bytes(), err
}

func (a *RsaCipher) EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

func (a *RsaCipher) DecodeURIComponent(str string) string {
	r, err := url.QueryUnescape(str)
	if err != nil {
		return str
	}
	return r
}
