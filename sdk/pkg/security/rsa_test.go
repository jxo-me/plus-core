package security

import (
	"testing"
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
)

func BenchmarkRsaEncryptAndDecrypt(b *testing.B) {
	plaintext := "这是测试aes-ctr加密内容256.#$5.*x,time:1653634619"
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
		if c.DecodeURIComponent(string(decrypt)) != plaintext {
			b.Error(`RSA Decrypt fail`)
		}
	}
}
