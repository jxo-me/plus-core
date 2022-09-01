package tus

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
)

func (h *sTus) log(ctx context.Context, eventName string, details ...string) {
	LogEvent(ctx, h.logger, eventName, details...)
}

func LogEvent(ctx context.Context, logger *glog.Logger, eventName string, details ...string) {
	result := make([]byte, 0, 100)

	result = append(result, `event="`...)
	result = append(result, eventName...)
	result = append(result, `" `...)

	for i := 0; i < len(details); i += 2 {
		result = append(result, details[i]...)
		result = append(result, `="`...)
		result = append(result, details[i+1]...)
		result = append(result, `" `...)
	}

	result = append(result, "\n"...)
	logger.Info(ctx, string(result))
}
