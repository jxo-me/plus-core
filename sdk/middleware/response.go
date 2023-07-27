package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/plus-core/pkg/response"
	"net/http"
)

// Response is the default middleware handling handler response object and its error.
func Response(r *ghttp.Request) {
	r.Middleware.Next()
	//glog.Warning(r.GetCtx(), "Response start ......")
	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}
	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)
	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		response.JsonExit(r, code.Code(), err.Error())
	} else if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		msg = http.StatusText(r.Response.Status)
		switch r.Response.Status {
		case http.StatusNotFound:
			code = gcode.CodeNotFound
		case http.StatusForbidden:
			code = gcode.CodeNotAuthorized
		default:
			code = gcode.CodeUnknown
		}
	} else {
		code = gcode.CodeOK
		msg = "ok"
	}
	//glog.Warning(r.GetCtx(), "Response end ......")
	response.JsonExit(r, code.Code(), msg, res)
}
