package registry

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type ServerRegistry struct {
	registry[*ghttp.Server]
}
