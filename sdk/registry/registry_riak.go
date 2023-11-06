package registry

import (
	"github.com/zegl/goriak/v3"
)

type RiakRegistry struct {
	registry[*goriak.Session]
}
