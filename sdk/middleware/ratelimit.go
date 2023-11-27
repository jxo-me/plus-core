package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/plus-core/pkg/v2/ratelimit/ip"
	"github.com/jxo-me/plus-core/pkg/v2/response"
	"github.com/jxo-me/plus-core/sdk/v2/errors/code"
)

func IpRateLimit(opts ...ip.Option) ghttp.HandlerFunc {
	rateLimit := ip.NewRateLimiter(opts...)
	return func(r *ghttp.Request) {
		// r.RemoteAddr 127.0.0.1:53274 r.GetClientIp 127.0.0.1
		var clientIp = r.GetClientIp()
		for _, excludeIp := range rateLimit.GetExcludes() {
			if clientIp == excludeIp {
				r.Middleware.Next()
				return
			}
		}

		limiter := rateLimit.Get(clientIp)
		if !limiter.Allow() {
			// rejected
			e := code.CodeLimitExceed
			r.Response.Status = e.Code()
			r.Response.WriteJsonExit(response.JsonRes{Code: e.Code(), Message: e.Message(), Data: g.Map{}})
		}
		// allowed
		r.Middleware.Next()
	}
}
