package ip

import (
	"github.com/jxo-me/plus-core/pkg/v2"
	"golang.org/x/time/rate"
	"sync"
)

// options of ip limiter.
type options struct {
	bucket int `json:"bucket"`
	// refilled at rate “r” tokens per second.
	limit rate.Limit `json:"limit"`
	// It implements a “token bucket” of size “b”
	burst int `json:"burst"`
	// exclude ip list
	exclude []string `json:"exclude"`
}

// WithSecondLimit with Limit is represented as number of events per second.
func WithSecondLimit(perSecond rate.Limit) Option {
	return func(o *options) {
		o.limit = perSecond
	}
}

// WithMaxQuota bursts of at most b tokens.
func WithMaxQuota(max int) Option {
	return func(o *options) {
		o.burst = max
	}
}

// WithExcludes bursts of at most b tokens.
func WithExcludes(ips []string) Option {
	return func(o *options) {
		o.exclude = ips
	}
}

// WithBucket max bucket num.
func WithBucket(size int) Option {
	return func(o *options) {
		o.bucket = size
	}
}

// Option is IpRateLimit option.
type Option func(*options)

// RateLimiter  ip rate limiter
type RateLimiter struct {
	ips  []*pkg.Bucket[*rate.Limiter]
	opts options
}

func NewRateLimiter(opts ...Option) *RateLimiter {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}
	if opt.bucket == 0 {
		opt.bucket = 1
	}
	buckets := make([]*pkg.Bucket[*rate.Limiter], opt.bucket)
	for i := 0; i < opt.bucket; i++ {
		buckets[i] = &pkg.Bucket[*rate.Limiter]{
			Mu:  &sync.RWMutex{},
			Idx: i,
			M:   make(map[string]*rate.Limiter),
		}
	}
	i := &RateLimiter{
		ips:  buckets,
		opts: opt,
	}

	return i
}

// Set creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *RateLimiter) Set(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(i.opts.limit, i.opts.burst)
	bucket := i.getBucket(ip)
	bucket.Set(ip, limiter)

	return limiter
}

// Get returns the rate limiter for the provided IP address if it exists.
// Otherwise, calls AddIP to add IP address to the map
func (i *RateLimiter) Get(ip string) *rate.Limiter {
	bucket := i.getBucket(ip)
	if !bucket.Has(ip) {
		limiter := rate.NewLimiter(i.opts.limit, i.opts.burst)
		bucket.Set(ip, limiter)
	}

	return bucket.Get(ip)
}

func (i *RateLimiter) GetExcludes() []string {
	return i.opts.exclude
}

// getBucket 根据用户名计算哈希值并映射到桶索引
func (i *RateLimiter) getBucket(key string) *pkg.Bucket[*rate.Limiter] {
	index := i.hashIndex(key)
	return i.ips[index]
}

func (i *RateLimiter) hashIndex(key string) int {
	hash := 0
	for _, ch := range key {
		hash = int(ch)
	}

	return hash % len(i.ips)
}
