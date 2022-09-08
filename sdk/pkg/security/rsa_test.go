package security

import (
	"fmt"
	"testing"
	"time"
)

const (
	ServerRsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDmzvnyyKiUcTkMpJIzW7e+svMq
a8KEu8zoq8iPbb3SUtDbZjEgbq+yUGUhxRi84HWD72dmBZfQM6xbewf3sirb5nMw
YcQcVRqSd7pqPIYbslLNFIE+l97p//egzEH0TdJSlEKA85JY8SKemznmc6OLXBjl
GC6GZlh1Hc/PnXsgvQIDAQAB
-----END PUBLIC KEY-----`
	ServerRsaPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDmzvnyyKiUcTkMpJIzW7e+svMqa8KEu8zoq8iPbb3SUtDbZjEg
bq+yUGUhxRi84HWD72dmBZfQM6xbewf3sirb5nMwYcQcVRqSd7pqPIYbslLNFIE+
l97p//egzEH0TdJSlEKA85JY8SKemznmc6OLXBjlGC6GZlh1Hc/PnXsgvQIDAQAB
AoGAKm8pMwLDQ69hTrq/CmQ1bwEVhdoUBDVG4gwbLot9l7qUHZI3PAA9myn+veuv
taIK2Pvd5brng2bMrHD9MQP3rubUUzmtzvWZSHk039MgVQp4WKnrKWqKNxquNAvx
xzRP+MOxYKmT83CJUN681kQHVjVuMPlaELfQm1ALofFNiUECQQD3/+MSl0nrsIjx
e/1IMGUeo+NbL+OVFPNJAFExEVHJ75aD5XWgk4DtykMT4C/K1aJTeJTPOBw7OSec
Wso45FhxAkEA7kEdzUXWOj3dMZGgZHMdah3hidBslpESrWt4x2hnZ8BglVJOomGD
NSb+O94Z0ezGygZ61JKuit+qowciT5ZTDQJBAOQxiTBfEv6Str909k7JcRIrfRU3
0HIqD93ZM9wxco+cLGP67Cwv7Y9f8C7Gt1MtFis2Jb8ygw+ffKorBC4rNpECQQDQ
5ulxjhVbddjXWJ+w07przA1wX/6GOmeNBRmehm18bBKDeuqVuDhqR2lNJt2u0hQv
GWVjs8U04Q5l6aBs9hqVAkBKgLlxPL9liwspWwSvwXH42iV1Qfuq0HxEkme8qMMh
Y6JATT3+fwHdfIQ/ZLOtUJnDW0Y0Z/ScroYAwuBK20TF
-----END RSA PRIVATE KEY-----`
	ClientRsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDIswTRdjsY315O3o3XlLaI2FCL
Usy8YlZ164pEtXzMxv756+FIX1ffGBJ1c4JXp4gar2INAmvAYEJi0ELQo4ilJl/y
4n1IWH5T8DwH8qlZeOzxL7CwmtkbE28RlIepMAfPpTjWw3XrK+IP0g1y+F6di/Cy
5QAL+iWDLh+GTQERpQIDAQAB
-----END PUBLIC KEY-----`
	ClientRsaPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDIswTRdjsY315O3o3XlLaI2FCLUsy8YlZ164pEtXzMxv756+FI
X1ffGBJ1c4JXp4gar2INAmvAYEJi0ELQo4ilJl/y4n1IWH5T8DwH8qlZeOzxL7Cw
mtkbE28RlIepMAfPpTjWw3XrK+IP0g1y+F6di/Cy5QAL+iWDLh+GTQERpQIDAQAB
AoGABEf9VFKHVDJ/moMY135unmCu5ynvAB3A5mcN3gVZEi00hzRG6/pMr4+d5S9/
mksSli8jkk946POktuLmafjzEzke+Vz1f64c7E83tXw1dK1/iLN6s8OdHn/aIR8I
j9ddfKscbA4+pXZSG0VUL4Lex0nkuaSfWfNRtrYXwgMbloECQQDrQMKEZ08b8XPB
KPNofLinQmWI0IZfwQwhqTbjyqiuskwr79avDkrmQwdCTR0I/kkLEWcrNArbg68T
J5tceIcxAkEA2mYmfcppYFTDl3Ln37ObhFnv0LxpiDVWLtHDCcrCQhkUBy+EK1ll
fzMRdGIPDtHVqYWnCdIgIIZa6ooG+zo8tQJAeSOQxoM4Hkh39cgzhuNqZl5vUEMo
XphWueKbrDLIJ6USSjGnV99BPE7Zpw90WxQt4rAkcv/Kc3zcoz1q5/P8kQJAWAv4
lpvksa5akBMGAfyRxODROtDKptwB26w7OhZhDG650U/JtAQ9U/ONpOLneS0FO2uj
dOUeiJ5Bxu5QdDX9RQJBALkJaLlC74EvcP5G89qN/eMSnuPF7ngVtLdhs/lyAj89
pPAbMC3Gv/yAKYjgI6dZ6g3x9Y8UGd1eUxZgRlmfVVI=
-----END RSA PRIVATE KEY-----`
)

func TestRsaEncryptAndDecrypts(t *testing.T) {
	c, err := NewRsaCiphers(
		ClientRsaPublicKey, ServerRsaPrivateKey,
		ServerRsaPublicKey, ClientRsaPrivateKey,
	)
	if err != nil {
		t.Error("NewRsaCiphers err", err)
	}
	nonce := time.Now().Unix()
	var tests = []struct {
		plaintext string
	}{
		{fmt.Sprintf("这是测试RSA加密内容.#$5.*x,time:%d", nonce)},
	}
	for _, test := range tests {
		// client Encrypt
		encrypt, err := c.Encrypt(test.plaintext)
		if err != nil {
			t.Error("RSA RsaCiphers Encrypt err", err)
		}
		// server Decrypt
		decrypt, err := c.Decrypt(encrypt)
		if err != nil {
			t.Error("RSA RsaCiphers Decrypt err", err)
		}
		if string(decrypt) != test.plaintext {
			t.Error(`RSA RsaCiphers Decrypt fail`)
		}
	}
}

func TestRsaEncryptAndDecrypt(t *testing.T) {
	c, err := NewRsaCiphers(
		ServerRsaPublicKey, ServerRsaPrivateKey,
		ClientRsaPublicKey, ClientRsaPrivateKey,
	)
	if err != nil {
		t.Error("NewRsaCiphers err", err)
	}
	nonce := time.Now().Unix()
	var tests = []struct {
		plaintext string
	}{
		{fmt.Sprintf("这是测试RSA加密内容.#$5.*x,time:%d", nonce)},
	}
	for _, test := range tests {
		// client test
		encrypt, err := c.ClientCipher.Encrypt(test.plaintext)
		if err != nil {
			t.Error("RSA ClientCipher Encrypt err", err)
		}
		decrypt, err := c.ClientCipher.Decrypt(encrypt)
		if err != nil {
			t.Error("RSA ClientCipher Decrypt err", err)
		}
		if string(decrypt) != test.plaintext {
			t.Error(`RSA ClientCipher Decrypt fail`)
		}
		// server test
		encrypt, err = c.ServerCipher.Encrypt(test.plaintext)
		if err != nil {
			t.Error("RSA ServerCipher Encrypt err", err)
		}
		decrypt, err = c.ServerCipher.Decrypt(encrypt)
		if err != nil {
			t.Error("RSA ServerCipher Decrypt err", err)
		}
		if string(decrypt) != test.plaintext {
			t.Error(`RSA ServerCipher Decrypt fail`)
		}
	}
}

func BenchmarkRsaEncryptAndDecrypt(b *testing.B) {
	plaintext := "这是测试RSA加密内容256.#$5.*x,time:1653634619"
	c, err := NewRsaCipher(ServerRsaPublicKey, ServerRsaPrivateKey)
	if err != nil {
		b.Error("RSA NewRsaCipher err", err)
	}
	for i := 0; i < b.N; i++ {
		ciphertext, err := c.Encrypt(plaintext)
		if err != nil {
			b.Error("RSA Encrypt err", err)
		}
		decrypt, err := c.Decrypt(ciphertext)
		if err != nil {
			b.Error("RSA Decrypt err", err)
		}
		if string(decrypt) != plaintext {
			b.Error(`RSA Decrypt fail`)
		}
	}
}
