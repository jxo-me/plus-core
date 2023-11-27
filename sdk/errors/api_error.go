package errors

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/jxo-me/plus-core/sdk/v2/errors/code"
)

func New(number int, message string) error {
	return gerror.NewCodeSkip(code.New(number, message, nil), code.MaxStackDepth)
}

func WithCode(co gcode.Code) error {
	return gerror.NewCodeSkip(code.WithCode(co, co.Detail()), code.MaxStackDepth)
}
