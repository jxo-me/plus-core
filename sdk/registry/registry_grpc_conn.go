package registry

import (
	"google.golang.org/grpc"
)

type GrpcConnRegistry struct {
	registry[*grpc.ClientConn]
}
