package tus

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
	"strconv"
)

func (h *Uploader) Middleware(r *ghttp.Request) {
	ctx := r.GetCtx()
	// Allow overriding the HTTP method. The reason for this is
	// that some libraries/environments to not support PATCH and
	// DELETE requests, e.g. Flash in a browser and parts of Java
	if newMethod := r.Header.Get("X-HTTP-Method-Override"); newMethod != "" {
		r.Method = newMethod
	}

	h.log(ctx, "RequestIncoming", "method", r.Method, "path", r.URL.Path, "requestId", getRequestId(r))

	h.Metrics.incRequestsTotal(r.Method)

	header := r.Response.Header()

	if origin := r.Header.Get("Origin"); origin != "" {
		header.Set("Access-Control-Allow-Origin", origin)

		if r.Method == "OPTIONS" {
			allowedMethods := "POST, HEAD, PATCH, OPTIONS"
			if !h.config.DisableDownload {
				allowedMethods += ", GET"
			}

			if !h.config.DisableTermination {
				allowedMethods += ", DELETE"
			}

			// Preflight request
			header.Add("Access-Control-Allow-Methods", allowedMethods)
			header.Add("Access-Control-Allow-Headers", "Authorization, Origin, X-Requested-With, X-Request-ID, X-HTTP-Method-Override, Content-Type, Upload-Length, Upload-Offset, Tus-Resumable, Upload-Metadata, Upload-Defer-Length, Upload-Concat")
			header.Set("Access-Control-Max-Age", "86400")

		} else {
			// Actual request
			header.Add("Access-Control-Expose-Headers", "Upload-Offset, Location, Upload-Length, Tus-Version, Tus-Resumable, Tus-Max-Size, Tus-Extension, Upload-Metadata, Upload-Defer-Length, Upload-Concat")
		}
	}

	// Set current version used by the server
	header.Set("Tus-Resumable", "1.0.0")

	// Add nosniff to all responses https://golang.org/src/net/http/server.go#L1429
	header.Set("X-Content-Type-Options", "nosniff")

	// Set appropriated headers in case of OPTIONS method allowing protocol
	// discovery and end with an 204 No Content
	if r.Method == "OPTIONS" {
		if h.config.MaxSize > 0 {
			header.Set("Tus-Max-Size", strconv.FormatInt(h.config.MaxSize, 10))
		}

		header.Set("Tus-Version", "1.0.0")
		header.Set("Tus-Extension", h.extensions)

		// Although the 204 No Content status code is a better fit in this case,
		// since we do not have a response body included, we cannot use it here
		// as some browsers only accept 200 OK as successful response to a
		// preflight request. If we send them the 204 No Content the response
		// will be ignored or interpreted as a rejection.
		// For example, the Presto engine, which is used in older versions of
		// Opera, Opera Mobile and Opera Mini, handles CORS this way.
		h.sendResp(r, http.StatusOK)
		return
	}

	// Test if the version sent by the client is supported
	// GET and HEAD methods are not checked since a browser may visit this URL and does
	// not include this header. GET requests are not part of the specification.
	if r.Method != "GET" && r.Method != "HEAD" && r.Header.Get("Tus-Resumable") != "1.0.0" {
		h.sendError(r, ErrUnsupportedVersion)
		return
	}
	r.Middleware.Next()
}
