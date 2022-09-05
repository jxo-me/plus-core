package security

import (
	"crypto/rc4"
	"encoding/base64"
	"encoding/hex"
	"math"
	"math/rand"
	"strings"
)

const salt uint8 = 67
const subLen int = 76

type Rc4Cipher struct {
	Key string
}

func NewRc4Cipher(key string) *Rc4Cipher {
	return &Rc4Cipher{Key: key}
}

func (r *Rc4Cipher) Encrypt(plaintext string) (string, error) {
	c, err := rc4.NewCipher([]byte(r.Key))
	if err != nil {
		return "", err
	}
	src := []byte(plaintext)
	dst := make([]byte, len(src))
	c.XORKeyStream(dst, src)
	return hex.EncodeToString(dst), nil
}

func (r *Rc4Cipher) Decrypt(ciphertext string) ([]byte, error) {
	src, err := hex.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	c, err := rc4.NewCipher([]byte(r.Key))
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(src))
	c.XORKeyStream(dst, src)

	//return string(dst), nil
	return dst, nil
}

// Rc4ClientEncrypt Client RC4 Config file Encrypt
func (r *Rc4Cipher) Rc4ClientEncrypt(plaintext string) (string, error) {
	bw := make([]byte, 0)
	pad := []uint8{2, 0, 0, 0}
	bw = append(bw, pad...)
	key := make([]uint8, 16)
	for i := 0; i < 16; i++ {
		c := uint8(r.RandInt(1, 256))
		s := (c ^ salt) % c
		bw = append(bw, c)
		key[i] = s
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return "", err
	}
	src := []byte(plaintext)
	dst := make([]byte, len(src))
	c.XORKeyStream(dst, src)
	bw = append(bw, dst...)
	return base64.StdEncoding.EncodeToString(bw), nil
}

// Rc4ClientDecrypt Client RC4 Config file Decrypt
func (r *Rc4Cipher) Rc4ClientDecrypt(ciphertext string) (string, error) {
	src, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	key := make([]uint8, 16)
	for i := 0; i < 16; i++ {
		s := src[i+4]
		key[i] = (s ^ salt) % s
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return "", err
	}
	dst := make([]byte, len(src)-20)
	c.XORKeyStream(dst, src[20:])

	return string(dst), nil
}

// CiphertextFormat Client encrypt content format
func (r *Rc4Cipher) CiphertextFormat(ciphertext string) string {
	str := ""
	l := int(math.Ceil(float64(len(ciphertext)) / float64(subLen)))
	for i := 0; i < l; i++ {
		if (i+1)*subLen < len(ciphertext) {
			str += ciphertext[i*subLen:(i+1)*subLen] + "\n"
		} else {
			str += ciphertext[i*subLen:]
		}
	}

	return str
}

func (r *Rc4Cipher) CiphertextReplace(ciphertext string) string {
	return strings.Replace(ciphertext, "\n", "", -1)
}

func (r *Rc4Cipher) RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}
