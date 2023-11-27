package code

import "github.com/gogf/gf/v2/errors/gcode"

const (
	// MaxStackDepth marks the max stack depth for error back traces.
	MaxStackDepth = 1
)

var (
	CodeLimitExceed = ApiCode{code: 429, message: "service unavailable due to rate limit exceeded", detail: nil} // is service unavailable due to rate limit exceeded.
)

// New creates and returns an error code.
// Note that it returns an interface object of Code.
func New(code int, message string, detail interface{}) gcode.Code {
	return ApiCode{
		code:    code,
		message: message,
		detail:  detail,
	}
}

// WithCode creates and returns a new error code based on given Code.
// The code and message is from given `code`, but the detail if from given `detail`.
func WithCode(code gcode.Code, detail interface{}) gcode.Code {
	return ApiCode{
		code:    code.Code(),
		message: code.Message(),
		detail:  detail,
	}
}
