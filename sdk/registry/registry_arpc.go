package registry

import (
	"github.com/lesismal/arpc"
)

type ArpcServerRegistry struct {
	registry[*arpc.Server]
}

type ArpcClientRegistry struct {
	registry[*arpc.Client]
}
