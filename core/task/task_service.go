package task

import (
	"context"
)

type IService interface {
	String() string
	Start(ctx context.Context)
}

type IServiceAppend interface {
	AppendStart(ctx context.Context, Routers []RabbitMqTask)
}
