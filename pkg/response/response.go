package response

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// JsonRes 数据返回通用JSON数据结构
type JsonRes struct {
	Code      int         `json:"code" dc:"Status Code" example:"0"` // 错误码((0:成功, 1:失败, >1:错误码))
	Message   string      `json:"message" dc:"Tips" example:"ok"`    // 提示信息
	Data      interface{} `json:"data" dc:"Data" example:"{}"`       // 返回数据(业务接口定义具体数据结构)
	RequestId string      `json:"request_id,omitempty" dc:"TraceId" example:"8a6f2805b4515ac12058e79e66539be9"`
}

// Json 返回标准JSON数据。
func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	} else {
		responseData = g.Map{}
	}
	r.Response.WriteJson(JsonRes{
		RequestId: gctx.CtxId(r.GetCtx()),
		Code:      code,
		Message:   message,
		Data:      responseData,
	})
}

// JsonExit 返回标准JSON数据并退出当前HTTP执行函数。
func JsonExit(r *ghttp.Request, code int, message string, data ...interface{}) {
	Json(r, code, message, data...)
	r.Exit()
}

// JsonRedirect 返回标准JSON数据引导客户端跳转。
func JsonRedirect(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(JsonRes{
		RequestId: gctx.CtxId(r.GetCtx()),
		Code:      code,
		Message:   message,
		Data:      responseData,
		Redirect:  redirect,
	})
}

// JsonRedirectExit 返回标准JSON数据引导客户端跳转，并退出当前HTTP执行函数。
func JsonRedirectExit(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	JsonRedirect(r, code, message, redirect, data...)
	r.Exit()
}
