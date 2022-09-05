package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jxo-me/plus-core/sdk/pkg/response"
	"github.com/jxo-me/plus-core/sdk/pkg/security"
	"io"
)

var insAesCrypto = AesCrypto{
	Key:   security.DefaultAesKey,
	NBits: security.DefaultAesNBits,
}

type AesCrypto struct {
	Cipher *security.AesCipher
	Key    string
	NBits  int
}

func Crypto(ctx context.Context) *AesCrypto {
	aes, err := security.NewAesCipher(insAesCrypto.Key, insAesCrypto.NBits)
	if err != nil {
		glog.Error(ctx, "NewCipher err", err)
	}
	insAesCrypto.Cipher = aes
	return &insAesCrypto
}

// DecryptRequest Aes解密请求数据
func (a *AesCrypto) DecryptRequest(r *ghttp.Request) {
	ctx := r.GetCtx()
	var (
		err error
	)
	buf, err := io.ReadAll(r.Request.Body)
	if err != nil {
		glog.Error(ctx, "Crypto io.ReadAll Request.Body err", err)
	}
	ciphertext := string(buf)
	glog.Debug(ctx, "raw request:", ciphertext)
	decrypt, err := a.Cipher.Decrypt(ciphertext)
	if err != nil {
		glog.Errorf(ctx, "RsaDecrypt error:%v", err)
	} else {
		//glog.Debug(ctx, "decrypt:", decrypt)
		r.Request.Body = io.NopCloser(bytes.NewReader(decrypt))
		r.Request.ContentLength = int64(len(decrypt))
	}
}

// EncryptResponse Aes加密返回数据
func (a *AesCrypto) EncryptResponse(r *ghttp.Request) {
	ctx := r.GetCtx()
	bf := r.Response.Buffer()
	glog.Debug(ctx, "raw Response body:", bf)
	// decode json response
	res := response.JsonRes{}
	err := json.Unmarshal(bf, &res)
	if err != nil {
		glog.Errorf(ctx, `json Unmarshal error: %+v`, err)
	}
	str := gconv.String(res.Data)
	encryptData, err := a.Cipher.Encrypt(str)
	if err != nil {
		glog.Warningf(ctx, `RsaEncrypt error: %+v`, err)
	} else {
		res.Data = encryptData
		glog.Debug(ctx, "aes.Encrypt Response body:", encryptData)
	}
	// Override response body
	r.Response.SetBuffer([]byte(gconv.String(res)))
}
