package security

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestEncryptAndDecrypt(t *testing.T) {
	rand.NewSource(time.Now().UnixNano())
	nonce := time.Now().Unix()
	var tests = []struct {
		key       string
		nBits     int
		plaintext string
	}{
		{"159054a86e3bfb85b5f1991cdb07645e", 256, fmt.Sprintf("这是测试aes-ctr加密内容256.#$5.*x,time:%d", nonce)},
		{"159054a86e3bfb85b5f1991cdb07645e", 192, fmt.Sprintf("这是测试aes-ctr加密内容192.#$5.*x,time:%d", nonce)},
		{"159054a86e3bfb85b5f1991cdb07645e", 128, fmt.Sprintf("这是测试aes-ctr加密内容128.#$5.*x,time:%d", nonce)},
	}
	for _, test := range tests {
		aes, err := NewAesCipher(test.key, test.nBits)
		if err != nil {
			t.Error("NewCipher err", err)
		}
		encode, err := aes.Encrypt(test.plaintext)
		if err != nil {
			t.Error("Encrypt err", err)
		}
		//return nil
		decode, err := aes.Decrypt(encode)
		if err != nil {
			t.Error("Decrypt err", err)
		}
		if string(decode) != test.plaintext {
			t.Error(`Decrypt fail`)
		}
	}
}

func TestDecrypt(t *testing.T) {
	var tests = []struct {
		key        string
		nBits      int
		ciphertext string
		text       string
	}{
		{"159054a86e3bfb85b5f1991cdb07645e", 256, "awI/NII7GQAu0s0w4gwN/1e4htQOkgYcCxUUdbxk7x7/nvwwq6ppLgdMDulefdPWYEgr/PhLD/SzA0WeCJgi8DM=", "这是测试aes-ctr加密内容256.#$5.*x,time:1653634619"},
		{"159054a86e3bfb85b5f1991cdb07645e", 192, "awI/NII7GQDuALyik9ECUmia3dAZA0bWD1PLqn6N+IguhdcsN1yvV0IFEDmEC/iMH14u2TG96/HrBD857dURvOg=", "这是测试aes-ctr加密内容192.#$5.*x,time:1653634619"},
		{"159054a86e3bfb85b5f1991cdb07645e", 128, "awI/NII7GQBp+P+xhXyWUkdrBbSVWRgMh1Ot8pfHrNU0ig+X3bH1SEMWMDiIerhwevvw6rV0VN/Q7Jv8lRtzbHE=", "这是测试aes-ctr加密内容128.#$5.*x,time:1653634619"},
	}
	for _, test := range tests {
		aes, err := NewAesCipher(test.key, test.nBits)
		if err != nil {
			t.Error("NewCipher err", err)
		}
		decode, err := aes.Decrypt(test.ciphertext)
		if err != nil {
			t.Error("Decrypt err", err)
		}
		if string(decode) != test.text {
			t.Error(`Decrypt fail`, test.nBits, string(decode))
		}
	}
}

func BenchmarkEncrypt256(b *testing.B) {
	key := "159054a86e3bfb85b5f1991cdb07645e"
	nBits := 256
	plaintext := "这是测试aes-ctr加密内容256.#$5.*x,time:1653634619"
	aes, err := NewAesCipher(key, nBits)
	if err != nil {
		b.Error("NewCipher err", err)
	}
	for i := 0; i < b.N; i++ {
		ciphertext, err := aes.Encrypt(plaintext)
		if err != nil {
			b.Error("Encrypt err", err)
		}
		_, err = aes.Decrypt(ciphertext)
		if err != nil {
			b.Error("Decrypt err", err)
		}
	}
}

func BenchmarkEncrypt192(b *testing.B) {
	key := "159054a86e3bfb85b5f1991cdb07645e"
	nBits := 192
	plaintext := "这是测试aes-ctr加密内容256.#$5.*x,time:1653634619"
	aes, err := NewAesCipher(key, nBits)
	if err != nil {
		b.Error("NewCipher err", err)
	}
	for i := 0; i < b.N; i++ {
		ciphertext, err := aes.Encrypt(plaintext)
		if err != nil {
			b.Error("Encrypt err", err)
		}
		_, err = aes.Decrypt(ciphertext)
		if err != nil {
			b.Error("Decrypt err", err)
		}
	}
}

func BenchmarkEncrypt128(b *testing.B) {
	key := "159054a86e3bfb85b5f1991cdb07645e"
	nBits := 128
	plaintext := "这是测试aes-ctr加密内容256.#$5.*x,time:1653634619"
	aes, err := NewAesCipher(key, nBits)
	if err != nil {
		b.Error("NewCipher err", err)
	}
	for i := 0; i < b.N; i++ {
		ciphertext, err := aes.Encrypt(plaintext)
		if err != nil {
			b.Error("Encrypt err", err)
		}
		_, err = aes.Decrypt(ciphertext)
		if err != nil {
			b.Error("Decrypt err", err)
		}
	}
}
