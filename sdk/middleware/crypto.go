package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jxo-me/plus-core/sdk/pkg/response"
	"github.com/jxo-me/plus-core/sdk/pkg/security"
	"io"
)

var insCrypto = crypto{}

type crypto struct {
	Cipher security.Crypto
}

func Crypto(cipher security.Crypto) *crypto {
	insCrypto.Cipher = cipher
	return &insCrypto
}

// DecryptRequest 解密请求数据
func (c *crypto) DecryptRequest(r *ghttp.Request) {
	ctx := r.GetCtx()
	var (
		err error
	)
	buf, err := io.ReadAll(r.Request.Body)
	if err != nil {
		glog.Error(ctx, "Crypto io.ReadAll Request.Body err", err)
	}
	ciphertext := string(buf)
	//glog.Debug(ctx, "raw request:", ciphertext)
	decrypt, err := c.Cipher.Decrypt(ciphertext)
	if err != nil {
		glog.Errorf(ctx, "RsaDecrypt error:%v", err)
	} else {
		//glog.Debug(ctx, "request Decrypt:", decrypt)
		r.Request.Body = io.NopCloser(bytes.NewReader(decrypt))
		r.Request.ContentLength = int64(len(decrypt))
	}
}

// EncryptResponse 加密响应数据
func (c *crypto) EncryptResponse(r *ghttp.Request) {
	ctx := r.GetCtx()
	bf := r.Response.Buffer()
	//glog.Debug(ctx, "raw Response body:", bf)
	// decode json response
	res := response.JsonRes{}
	err := json.Unmarshal(bf, &res)
	if err != nil {
		glog.Errorf(ctx, `json Unmarshal error: %+v`, err)
	}
	str := gconv.String(res.Data)
	encryptData, err := c.Cipher.Encrypt(str)
	if err != nil {
		glog.Warningf(ctx, `RsaEncrypt error: %+v`, err)
	} else {
		res.Data = encryptData
		//glog.Debug(ctx, "aes.Encrypt Response body:", encryptData)
	}
	// Override response body
	r.Response.SetBuffer([]byte(gconv.String(res)))
}
