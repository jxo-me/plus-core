package tus

import (
	"errors"
	"net/http"
)

var (
	ErrUnsupportedVersion               = NewHTTPError(errors.New("unsupported version"), http.StatusPreconditionFailed)
	ErrMaxSizeExceeded                  = NewHTTPError(errors.New("maximum size exceeded"), http.StatusRequestEntityTooLarge)
	ErrInvalidContentType               = NewHTTPError(errors.New("missing or invalid Content-Type header"), http.StatusBadRequest)
	ErrInvalidUploadLength              = NewHTTPError(errors.New("missing or invalid Upload-Length header"), http.StatusBadRequest)
	ErrInvalidOffset                    = NewHTTPError(errors.New("missing or invalid Upload-Offset header"), http.StatusBadRequest)
	ErrNotFound                         = NewHTTPError(errors.New("upload not found"), http.StatusNotFound)
	ErrFileLocked                       = NewHTTPError(errors.New("file currently locked"), 423) // Locked (WebDAV) (RFC 4918)
	ErrMismatchOffset                   = NewHTTPError(errors.New("mismatched offset"), http.StatusConflict)
	ErrSizeExceeded                     = NewHTTPError(errors.New("resource's size exceeded"), http.StatusRequestEntityTooLarge)
	ErrNotImplemented                   = NewHTTPError(errors.New("feature not implemented"), http.StatusNotImplemented)
	ErrUploadNotFinished                = NewHTTPError(errors.New("one of the partial uploads is not finished"), http.StatusBadRequest)
	ErrInvalidConcat                    = NewHTTPError(errors.New("invalid Upload-Concat header"), http.StatusBadRequest)
	ErrModifyFinal                      = NewHTTPError(errors.New("modifying a final upload is not allowed"), http.StatusForbidden)
	ErrUploadLengthAndUploadDeferLength = NewHTTPError(errors.New("provided both Upload-Length and Upload-Defer-Length"), http.StatusBadRequest)
	ErrInvalidUploadDeferLength         = NewHTTPError(errors.New("invalid Upload-Defer-Length header"), http.StatusBadRequest)
	ErrUploadStoppedByServer            = NewHTTPError(errors.New("upload has been stopped by server"), http.StatusBadRequest)

	errReadTimeout     = errors.New("read tcp: i/o timeout")
	errConnectionReset = errors.New("read tcp: connection reset by peer")
)

// HTTPError represents an error with an additional status code attached
// which may be used when this error is sent in an HTTP response.
// See the net/http package for standardized status codes.
type HTTPError interface {
	error
	StatusCode() int
	Body() []byte
}

type httpError struct {
	error
	statusCode int
}

func (err httpError) StatusCode() int {
	return err.statusCode
}

func (err httpError) Body() []byte {
	return []byte(err.Error())
}

// NewHTTPError adds the given status code to the provided error and returns
// the new error instance. The status code may be used in corresponding HTTP
// responses. See the net/http package for standardized status codes.
func NewHTTPError(err error, statusCode int) HTTPError {
	return httpError{err, statusCode}
}
