package middleware

import (
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/jxo-me/plus-core/sdk/pkg/security"
	"io/ioutil"
	"strings"
)

// AesDecryptRequest Aes解密请求数据
func AesDecryptRequest(r *ghttp.Request) {
	// data Decrypt request
	if r.Method == "POST" || r.Method == "PUT" {
		var (
			n   int
			err error
		)
		buf := make([]byte, r.Request.ContentLength)
		n, err = r.Request.Body.Read(buf)
		if err != nil {
			glog.Error(r.GetCtx(), "Request.Body.Read err", err)
		}
		ciphertext := string(buf[0:n])
		glog.Debug(r.GetCtx(), "raw:", ciphertext)
		aes, err := security.NewAesCipher(security.DefaultAesKey, security.DefaultAesNBits)
		if err != nil {
			glog.Error(r.GetCtx(), "NewCipher err", err)
		}
		decrypt, err := aes.Decrypt(ciphertext)
		//decrypt, err := utils.RsaDecrypt(ciphertext, consts.ClientRsaPrivateKey)
		if err != nil {
			glog.Errorf(r.GetCtx(), "RsaDecrypt error:%v", err)
		} else {
			// rewrite body
			newBodyContent := string(decrypt)
			r.Request.Body = ioutil.NopCloser(strings.NewReader(newBodyContent))
			r.Request.ContentLength = int64(len(newBodyContent))
		}

	}
	// Continue execution of next middleware.
	r.Middleware.Next()
}

// AesEncryptResponse Aes加密返回数据
func AesEncryptResponse(r *ghttp.Request) {
	r.Middleware.Next()

	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}

	// dataEncrypt response
	if r.Method == "POST" || r.Method == "PUT" {
		ctx := r.GetCtx()
		res := r.GetHandlerResponse()
		marshal, err := json.Marshal(res)
		if err != nil {
			glog.Warningf(ctx, `json Marshal error: %+v`, err)
		} else {
			aes, err := security.NewAesCipher(security.DefaultAesKey, security.DefaultAesNBits)
			if err != nil {
				glog.Error(ctx, "NewCipher err", err)
			}
			encryptData, err := aes.Encrypt(gconv.String(marshal))
			if err != nil {
				glog.Warningf(ctx, `RsaEncrypt error: %+v`, err)
			} else {
				// rewrite response
				newBodyContent := []byte(encryptData)
				r.Response.SetBuffer(newBodyContent)
			}
		}
	}

}
