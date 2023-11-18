package registry

import (
	"github.com/jxo-me/plus-core/pkg/v2/erlang"
)

type ErlangNodeRegistry struct {
	registry[*erlang.Node]
}
