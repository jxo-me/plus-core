package tus

import "github.com/gogf/gf/v2/net/ghttp"

// HookEvent represents an event from tus which can be handled by the application.
type HookEvent struct {
	// Upload contains information about the upload that caused this hook
	// to be fired.
	Upload FileInfo
	// HTTPRequest contains details about the HTTP request that reached
	// tus.
	HTTPRequest HTTPRequest
}

func newHookEvent(info FileInfo, r *ghttp.Request) HookEvent {
	return HookEvent{
		Upload: info,
		HTTPRequest: HTTPRequest{
			Method:     r.Method,
			URI:        r.RequestURI,
			RemoteAddr: r.RemoteAddr,
			Header:     r.Header,
		},
	}
}
