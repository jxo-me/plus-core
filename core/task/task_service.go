package task

import (
	"context"
)

type IService interface {
	String() string
	Start(ctx context.Context)
}
