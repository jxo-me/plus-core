package registry

import lockerLib "github.com/jxo-me/plus-core/core/locker"

type LockerRegistry struct {
	registry[lockerLib.ILocker]
}
