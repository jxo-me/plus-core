package tus

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Send the error in the response body. The status code will be looked up in
// ErrStatusCodes. If none is found 500 Internal Error will be used.
func (h *Uploader) sendError(r *ghttp.Request, err error) {
	// Errors for read timeouts contain too much information which is not
	// necessary for us and makes grouping for the metrics harder. The error
	// message looks like: read tcp 127.0.0.1:1080->127.0.0.1:53673: i/o timeout
	// Therefore, we use a common error message for all of them.
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		err = errReadTimeout
	}

	// Errors for connnection resets also contain TCP details, we don't need, e.g:
	// read tcp 127.0.0.1:1080->127.0.0.1:10023: read: connection reset by peer
	// Therefore, we also trim those down.
	if strings.HasSuffix(err.Error(), "read: connection reset by peer") {
		err = errConnectionReset
	}

	// TODO: Decide if we should handle this in here, in body_reader or not at all.
	// If the HTTP PATCH request gets interrupted in the middle (e.g. because
	// the user wants to pause the upload), Go's net/http returns an io.ErrUnexpectedEOF.
	// However, for the h it's not important whether the stream has ended
	// on purpose or accidentally.
	//if err == io.ErrUnexpectedEOF {
	//	err = nil
	//}

	// TODO: Decide if we want to ignore connection reset errors all together.
	// In some cases, the HTTP connection gets reset by the other peer. This is not
	// necessarily the tus client but can also be a proxy in front of tus, e.g. HAProxy 2
	// is known to reset the connection to tus, when the tus client closes the connection.
	// To avoid erroring out in this case and loosing the uploaded data, we can ignore
	// the error here without causing harm.
	//if strings.Contains(err.Error(), "read: connection reset by peer") {
	//	err = nil
	//}

	statusErr, ok := err.(HTTPError)
	if !ok {
		statusErr = NewHTTPError(err, http.StatusInternalServerError)
	}

	reason := append(statusErr.Body(), '\n')
	if r.Method == "HEAD" {
		reason = nil
	}
	r.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.Response.Header().Set("Content-Length", strconv.Itoa(len(reason)))
	r.Response.WriteStatus(statusErr.StatusCode(), reason)

	h.log(r.GetCtx(), "ResponseOutgoing", "status", strconv.Itoa(statusErr.StatusCode()), "method", r.Method, "path", r.URL.Path, "error", err.Error(), "requestId", getRequestId(r))
	glog.Warning(r.GetCtx(), "sendError:", reason)
	h.Metrics.incErrorsTotal(statusErr)
}

// sendResp writes the header to w with the specified status code.
func (h *Uploader) sendResp(r *ghttp.Request, status int, content ...interface{}) {
	r.Response.WriteStatus(status, content)
	h.log(r.GetCtx(), "ResponseOutgoing", "status", strconv.Itoa(status), "method", r.Method, "path", r.URL.Path, "requestId", getRequestId(r))
}
