package ws

import (
	"context"
)

type SrvBefore interface {
	Before(ctx context.Context, message *MessageReq) error
}

type SrvHandler interface {
	Handle(ctx context.Context, message *MessageReq) (*MessageRes, error)
}

type SrvAfter interface {
	After(ctx context.Context, response *MessageRes) error
}

type SrvAction interface {
	Action() Actions
}

type Service interface {
	SrvAction
}

type ActPost interface {
	Post(ctx context.Context, message *MessageReq) (*MessageRes, error)
}

type ActDelete interface {
	Delete(ctx context.Context, message *MessageReq) (*MessageRes, error)
}

type ActUpdate interface {
	Update(ctx context.Context, message *MessageReq) (*MessageRes, error)
}

type ActGet interface {
	Get(ctx context.Context, message *MessageReq) (*MessageRes, error)
}

type ActJoin interface {
	Join(ctx context.Context, message *MessageReq) (*MessageRes, error)
}

type ActLeave interface {
	Leave(ctx context.Context, message *MessageReq) (*MessageRes, error)
}

type Actions interface {
	ActPost
	ActDelete
	ActUpdate
	ActGet
	ActJoin
	ActLeave
}
