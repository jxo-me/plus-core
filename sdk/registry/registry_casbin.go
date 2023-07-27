package registry

import "github.com/casbin/casbin/v2"

type CasBinRegistry struct {
	registry[*casbin.SyncedEnforcer]
}
