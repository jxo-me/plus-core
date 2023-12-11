package send

import "context"

type ISendMsg interface {
	Format(level string) string
}

type ISender[T any] interface {
	String() string
	Info(ctx context.Context, msg T) (err error)
	Warn(ctx context.Context, msg T) (err error)
	Error(ctx context.Context, msg T) (err error)
}
