package tus

import "net/http"

// HTTPRequest contains basic details of an incoming HTTP request.
type HTTPRequest struct {
	// Method is the HTTP method, e.g. POST or PATCH
	Method string
	// URI is the full HTTP request URI, e.g. /files/fooo
	URI string
	// RemoteAddr contains the network address that sent the request
	RemoteAddr string
	// Header contains all HTTP headers as present in the HTTP request.
	Header http.Header
}
