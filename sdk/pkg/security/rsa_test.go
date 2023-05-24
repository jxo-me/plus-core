package security

import (
	"fmt"
	"testing"
	"time"
)

const (
	ServerRsaPublicKey  = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDmzvnyyKiUcTkMpJIzW7e+svMqa8KEu8zoq8iPbb3SUtDbZjEgbq+yUGUhxRi84HWD72dmBZfQM6xbewf3sirb5nMwYcQcVRqSd7pqPIYbslLNFIE+l97p//egzEH0TdJSlEKA85JY8SKemznmc6OLXBjlGC6GZlh1Hc/PnXsgvQIDAQAB`
	ServerRsaPrivateKey = `MIICXQIBAAKBgQDmzvnyyKiUcTkMpJIzW7e+svMqa8KEu8zoq8iPbb3SUtDbZjEgbq+yUGUhxRi84HWD72dmBZfQM6xbewf3sirb5nMwYcQcVRqSd7pqPIYbslLNFIE+l97p//egzEH0TdJSlEKA85JY8SKemznmc6OLXBjlGC6GZlh1Hc/PnXsgvQIDAQABAoGAKm8pMwLDQ69hTrq/CmQ1bwEVhdoUBDVG4gwbLot9l7qUHZI3PAA9myn+veuvtaIK2Pvd5brng2bMrHD9MQP3rubUUzmtzvWZSHk039MgVQp4WKnrKWqKNxquNAvxxzRP+MOxYKmT83CJUN681kQHVjVuMPlaELfQm1ALofFNiUECQQD3/+MSl0nrsIjxe/1IMGUeo+NbL+OVFPNJAFExEVHJ75aD5XWgk4DtykMT4C/K1aJTeJTPOBw7OSecWso45FhxAkEA7kEdzUXWOj3dMZGgZHMdah3hidBslpESrWt4x2hnZ8BglVJOomGDNSb+O94Z0ezGygZ61JKuit+qowciT5ZTDQJBAOQxiTBfEv6Str909k7JcRIrfRU30HIqD93ZM9wxco+cLGP67Cwv7Y9f8C7Gt1MtFis2Jb8ygw+ffKorBC4rNpECQQDQ5ulxjhVbddjXWJ+w07przA1wX/6GOmeNBRmehm18bBKDeuqVuDhqR2lNJt2u0hQvGWVjs8U04Q5l6aBs9hqVAkBKgLlxPL9liwspWwSvwXH42iV1Qfuq0HxEkme8qMMhY6JATT3+fwHdfIQ/ZLOtUJnDW0Y0Z/ScroYAwuBK20TF`
	ClientRsaPublicKey  = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDIswTRdjsY315O3o3XlLaI2FCLUsy8YlZ164pEtXzMxv756+FIX1ffGBJ1c4JXp4gar2INAmvAYEJi0ELQo4ilJl/y4n1IWH5T8DwH8qlZeOzxL7CwmtkbE28RlIepMAfPpTjWw3XrK+IP0g1y+F6di/Cy5QAL+iWDLh+GTQERpQIDAQAB`
	ClientRsaPrivateKey = `MIICXAIBAAKBgQDIswTRdjsY315O3o3XlLaI2FCLUsy8YlZ164pEtXzMxv756+FIX1ffGBJ1c4JXp4gar2INAmvAYEJi0ELQo4ilJl/y4n1IWH5T8DwH8qlZeOzxL7CwmtkbE28RlIepMAfPpTjWw3XrK+IP0g1y+F6di/Cy5QAL+iWDLh+GTQERpQIDAQABAoGABEf9VFKHVDJ/moMY135unmCu5ynvAB3A5mcN3gVZEi00hzRG6/pMr4+d5S9/mksSli8jkk946POktuLmafjzEzke+Vz1f64c7E83tXw1dK1/iLN6s8OdHn/aIR8Ij9ddfKscbA4+pXZSG0VUL4Lex0nkuaSfWfNRtrYXwgMbloECQQDrQMKEZ08b8XPBKPNofLinQmWI0IZfwQwhqTbjyqiuskwr79avDkrmQwdCTR0I/kkLEWcrNArbg68TJ5tceIcxAkEA2mYmfcppYFTDl3Ln37ObhFnv0LxpiDVWLtHDCcrCQhkUBy+EK1llfzMRdGIPDtHVqYWnCdIgIIZa6ooG+zo8tQJAeSOQxoM4Hkh39cgzhuNqZl5vUEMoXphWueKbrDLIJ6USSjGnV99BPE7Zpw90WxQt4rAkcv/Kc3zcoz1q5/P8kQJAWAv4lpvksa5akBMGAfyRxODROtDKptwB26w7OhZhDG650U/JtAQ9U/ONpOLneS0FO2ujdOUeiJ5Bxu5QdDX9RQJBALkJaLlC74EvcP5G89qN/eMSnuPF7ngVtLdhs/lyAj89pPAbMC3Gv/yAKYjgI6dZ6g3x9Y8UGd1eUxZgRlmfVVI=`
)

func TestRsaEncryptAndDecrypt(t *testing.T) {
	c, err := NewRsaCipher(ServerRsaPublicKey, ClientRsaPrivateKey)
	if err != nil {
		t.Error("Client NewRsaCiphers err", err)
	}
	s, err := NewRsaCipher(ClientRsaPublicKey, ServerRsaPrivateKey)
	if err != nil {
		t.Error("Server NewRsaCiphers err", err)
	}
	nonce := time.Now().Unix()
	var tests = []struct {
		plaintext string
	}{
		{fmt.Sprintf("这是测试RSA加密内容.#$5.*x,time:%d", nonce)},
	}
	for _, test := range tests {
		// client test
		encrypt, err := c.Encrypt(test.plaintext)
		if err != nil {
			t.Error("Client RSA ClientCipher Encrypt err", err)
		}
		decrypt, err := s.Decrypt(encrypt)
		if err != nil {
			t.Error("Server RSA ClientCipher Decrypt err", err)
		}
		if string(decrypt) != test.plaintext {
			t.Error(`RSA ClientCipher Decrypt fail`)
		}
		// server test
		encrypt, err = s.Encrypt(test.plaintext)
		if err != nil {
			t.Error("Server RSA ServerCipher Encrypt err", err)
		}
		decrypt, err = c.Decrypt(encrypt)
		if err != nil {
			t.Error("Client RSA ServerCipher Decrypt err", err)
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
