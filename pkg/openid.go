package pkg

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type OpenId struct {
	Key  string `json:"key"`   // 密匙
	KeyA string `json:"key_a"` // 密匙a会参与加解密
	KeyB string `json:"key_b"` // 密匙b会用来做数据完整性验证
	//KeyC   string `json:"key_c"`  // 密匙c用于变化生成的密文
	Salt   string `json:"salt"`   // salt
	Expiry int    `json:"expiry"` // 过期时间
	KeyLen int    `json:"key_len"`
}

func NewOpenId(keyA, keyB, salt string) *OpenId {
	return &OpenId{
		KeyA:   keyA,
		KeyB:   keyB,
		Salt:   salt,
		KeyLen: 4,
	}
}

func NewOpenIdWithKey(key string) *OpenId {
	keyHash := md5.Sum([]byte(key))
	keyA := fmt.Sprintf("%x", md5.Sum(keyHash[:16]))
	keyB := fmt.Sprintf("%x", md5.Sum(keyHash[16:]))
	//salt := fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))
	salt := fmt.Sprintf("%x", keyHash)
	return &OpenId{
		Key:    key,
		KeyA:   keyA,
		KeyB:   keyB,
		Salt:   salt,
		KeyLen: 4,
	}
}

func (s *OpenId) Encode(userId int) (string, error) {
	authcode, err := s.encode(fmt.Sprintf("%d", userId))
	if err != nil {
		return "", err
	}
	encoded := s.encodeURI(authcode)
	return strings.ReplaceAll(encoded, "%", "_"), nil
}

func (s *OpenId) encode(str string) (string, error) {
	//keyLen := 4
	//keya := "a6d7ad72c18ac0e15e5cd994b3cc79b4"
	//keyb := "5df37d77e57300cfb6992fd1857077a4"
	//times := "e4b38aed46581c30418803a63abb332a"

	//keyHash := md5.Sum([]byte(key))
	//keya = fmt.Sprintf("%x", md5.Sum(keyHash[:16]))
	//keyb = fmt.Sprintf("%x", md5.Sum(keyHash[16:]))
	//keyc := ""
	keyC := s.substr(s.Salt, len(s.Salt)-s.KeyLen, s.KeyLen)
	//keyc = fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))[32-s.KeyLen:]
	cryptKey := s.KeyA + s.md5Str(s.KeyA+keyC)
	cryptKeyLen := len(cryptKey)
	expiryStr := ""
	if s.Expiry == 0 {
		expiryStr = s.sprintf(s.Expiry)
	} else {
		expiryStr = s.sprintf(s.Expiry + int(time.Now().Unix()))
	}
	str1 := expiryStr + s.substr(s.md5Str(str+s.KeyB), 0, 16) + str
	str1Len := len(str1)
	result := s.core(cryptKey, cryptKeyLen, str1, str1Len)
	return keyC + strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(result)), "=", ""), nil
}

func (s *OpenId) Decode(OpenID string) (string, error) {
	OpenID = strings.ReplaceAll(OpenID, "_", "%")
	decoded, err := s.decodeURI(OpenID)
	if err != nil {
		return "", err
	}
	return s.decode(decoded)
}

func (s *OpenId) decode(str string) (string, error) {
	//keyLen := 4
	//keya := "a6d7ad72c18ac0e15e5cd994b3cc79b4"
	//keyb := "5df37d77e57300cfb6992fd1857077a4"
	//times := "e4b38aed46581c30418803a63abb332a"

	//keyHash := md5.Sum([]byte(key))
	//keya = fmt.Sprintf("%x", md5.Sum(keyHash[:16]))
	//keyb = fmt.Sprintf("%x", md5.Sum(keyHash[16:]))
	//keyC := ""
	keyC := str[:s.KeyLen]
	cryptKey := s.KeyA + s.md5Str(s.KeyA+keyC)
	cryptKeyLen := len(cryptKey)
	str1 := ""
	base, err := s.base64Decode(str[4:])
	if err != nil {
		return "", err
	}
	str1 = string(base)
	str1Len := len(str1)
	result := s.core(cryptKey, cryptKeyLen, str1, str1Len)
	return result[26:], nil
}

func (s *OpenId) authCode(str, operation, key string, expiry int) (string, error) {
	keyLen := 4
	keya := "a6d7ad72c18ac0e15e5cd994b3cc79b4"
	keyb := "5df37d77e57300cfb6992fd1857077a4"
	times := "e4b38aed46581c30418803a63abb332a"

	//keyHash := md5.Sum([]byte(key))
	//keya = fmt.Sprintf("%x", md5.Sum(keyHash[:16]))
	//keyb = fmt.Sprintf("%x", md5.Sum(keyHash[16:]))
	keyc := ""
	if operation == "DECODE" {
		keyc = str[:keyLen]
	} else {
		keyc = s.substr(times, len(times)-keyLen, keyLen)
		//keyc = fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))[32-keyLen:]
	}
	fmt.Println("keyc:", keyc)
	cryptKey := keya + s.md5Str(keya+keyc)
	fmt.Println("cryptKey:", cryptKey)
	cryptKeyLen := len(cryptKey)
	str1 := ""
	if operation == "DECODE" {
		base, err := s.base64Decode(str[4:])
		if err != nil {
			return "", err
		}
		str1 = string(base)
	} else {
		expiryStr := ""
		if expiry == 0 {
			expiryStr = s.sprintf(expiry)
		} else {
			expiryStr = s.sprintf(expiry + int(time.Now().Unix()))
		}
		str1 = expiryStr + s.substr(s.md5Str(str+keyb), 0, 16) + str
	}
	fmt.Println("str1:", str1)
	fmt.Println("str1 len:", len(str1))
	str1Len := len(str1)
	result := s.core(cryptKey, cryptKeyLen, str1, str1Len)
	//fmt.Println("result:", result)
	if operation == "DECODE" {
		return result[26:], nil
	} else {
		return keyc + strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(result)), "=", ""), nil
	}
}

func (s *OpenId) core(cryptKey string, cryptKeyLen int, str1 string, str1Len int) string {
	boxs := make([]int, 256)
	for i := 0; i < 256; i++ {
		boxs[i] = i
	}
	randKeys := s.genRandKey([]byte(cryptKey), cryptKeyLen)
	randBox := s.randBoxs(0, 0, 256, boxs, randKeys)

	result := s.coreEncode(0, 0, 0, []byte(str1), str1Len, 256, randBox)
	return string(result)
}

func (s *OpenId) genRandKey(cryptKey []byte, keyLen int) []int {
	randKey := make([]int, 256)
	for i := 0; i < 256; i++ {
		randKey[i] = int(cryptKey[i%keyLen])
	}
	return randKey
}

func (s *OpenId) randBoxs(j, i, max int, boxs, randKeys []int) []int {
	if i == max {
		return boxs
	}
	boxsi := boxs[i]
	randKeysi := randKeys[i]
	jj := (j + boxsi + randKeysi) % max
	tmp := boxsi
	boxs[i] = boxs[jj]
	boxs[jj] = tmp
	return s.randBoxs(jj, i+1, max, boxs, randKeys)
}

func (s *OpenId) encodeURI(input string) string {
	encodedURI := url.QueryEscape(input)
	return encodedURI
}

func (s *OpenId) decodeURI(encodedURI string) (string, error) {
	decodedURI, err := url.QueryUnescape(encodedURI)
	if err != nil {
		return "", err
	}
	return decodedURI, nil
}

func (s *OpenId) md5Str(str string) string {
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}

func (s *OpenId) sprintf(Expiry int) string {
	return fmt.Sprintf("%010d", Expiry)
}

func (s *OpenId) substr(str string, offset int, length int) string {
	if offset < 0 || offset >= len(str) || length <= 0 {
		return ""
	}

	end := offset + length
	if end > len(str) {
		end = len(str)
	}

	return str[offset:end]
}

func (s *OpenId) coreEncode(a, j, i int, str []byte, strLen, max int, boxs []int) []byte {
	if i == strLen {
		return str
	}
	aa := (a + 1) % max
	boxa := boxs[aa]
	jj := (j + boxa) % max
	tmp := boxa
	boxs[aa] = boxs[jj]
	boxs[jj] = tmp
	boxaa := boxs[aa]
	boxj := boxs[jj]
	n := (boxaa + boxj) % max
	boxn := boxs[n]
	str[i] = str[i] ^ byte(boxn)
	return s.coreEncode(aa, jj, i+1, str, strLen, max, boxs)
}

func (s *OpenId) base64Decode(encode string) ([]byte, error) {
	l := len(encode)
	m := l % 4
	padding := strings.Repeat("=", 4-m)

	if m == 0 {
		return base64.StdEncoding.DecodeString(encode)
	}
	// Ensure the input string is a multiple of 4 by padding with '=' characters
	encoded := encode + padding

	return base64.StdEncoding.DecodeString(encoded)
}
