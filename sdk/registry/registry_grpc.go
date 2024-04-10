package registry

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type GrpcRegistry struct {
	registry[*grpcx.GrpcServer]
}
