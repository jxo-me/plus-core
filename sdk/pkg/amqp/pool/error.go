package pool

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"runtime/debug"
)

// FailOnError log stack and fatal with given message
func FailOnError(ctx context.Context, err error, msg string) {
	if err != nil {
		stack := debug.Stack()
		g.Log().Fatalf(ctx, "%s: %s - stack:\n%s", msg, err, stack)
	}
}

// WrapError wrap a error with given message
func WrapError(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("%s: %s", msg, err)
	}
	return nil
}
