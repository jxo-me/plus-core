package otp

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/skip2/go-qrcode"
	"net/url"
	"strings"
	"time"
)

// _otpAuth is a one-time-password configuration.  This object will be modified by calls to
// Authenticate and should be saved to ensure the codes are in fact only used
// once. GoogleAuthenticator
type sOtpAuth struct {
	UTC       bool // use UTC for the timestamp instead of local time
	OTPIssuer string
	Logger    glog.ILogger
}

func New(issuer string, utc bool, log glog.ILogger) *sOtpAuth {
	return &sOtpAuth{
		UTC:       utc,
		OTPIssuer: issuer,
		Logger:    log,
	}
}

func (s *sOtpAuth) ProvisionURI(user string, secret string) string {
	return s.ProvisionURIWithIssuer(user, secret, s.OTPIssuer)
}

func (s *sOtpAuth) ProvisionURIWithIssuer(user string, secret string, issuer string) string {
	auth := "totp/"
	q := make(url.Values)
	q.Add("secret", secret)
	if issuer != "" {
		q.Add("issuer", issuer)
		auth += issuer + ":"
	}

	return "otpauth://" + auth + user + "?" + q.Encode()
}

// GetSecret
// 获取秘钥
func (s *sOtpAuth) GetSecret(ctx context.Context) string {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, s.un())
	if err != nil {
		s.Logger.Error(ctx, "binary write errors", err.Error())
	}
	return strings.ToUpper(s.base32encode(s.hmacSha1(buf.Bytes(), nil)))
}

func (s *sOtpAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (s *sOtpAuth) base32decode(str string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(str)
}

func (s *sOtpAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (s *sOtpAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (s *sOtpAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (s *sOtpAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := s.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := s.toUint32(hashParts)
	return number % 1000000
}

func (s *sOtpAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

// VerifyCode
// 验证动态码
func (s *sOtpAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := s.GetCode(secret)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}

// GetCode
// 获取动态码
func (s *sOtpAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := s.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	var t0 int64
	// assume we're on Time-based OTP
	if s.UTC {
		t0 = time.Now().UTC().Unix() / 30
	} else {
		t0 = time.Now().Unix() / 30
	}
	number := s.oneTimePassword(secretKey, s.toBytes(t0))
	return fmt.Sprintf("%06d", number), nil
}

func (s *sOtpAuth) GetQrcode(url string) (string, error) {
	var png []byte
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	str := s.imageToBase64(png)
	return str, nil
}

func (s *sOtpAuth) imageToBase64(img []byte) string {
	imgBase64Str := base64.StdEncoding.EncodeToString(img)
	return s.imgBase64Str(imgBase64Str)
}

func (s *sOtpAuth) imgBase64Str(base64Str string) string {
	return fmt.Sprintf("data:image/png;base64,%s", base64Str)
}
