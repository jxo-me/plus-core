package locker

import "github.com/jxo-me/redislock"

type ILocker interface {
	String() string
	Lock(key string, ttl int64, options ...redislock.Option) (*redislock.Mutex, error)
}
