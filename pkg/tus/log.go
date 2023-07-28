package tus

import (
	"context"
	"github.com/jxo-me/plus-core/core/v2/logger"
)

func (u *Uploader) log(ctx context.Context, eventName string, details ...string) {
	LogEvent(ctx, u.logger, eventName, details...)
}

func LogEvent(ctx context.Context, log logger.ILogger, eventName string, details ...string) {
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
	log.Infof(ctx, string(result))
}
