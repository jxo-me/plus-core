package registry

import lockerLib "github.com/jxo-me/plus-core/core/v2/locker"

type LockerRegistry struct {
	registry[lockerLib.ILocker]
}
