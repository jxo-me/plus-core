package ip

import (
	"golang.org/x/time/rate"
	"sync"
)

// options of ip limiter.
type options struct {
	// refilled at rate “r” tokens per second.
	limit rate.Limit
	// It implements a “token bucket” of size “b”
	burst int
	// exclude ip list
	Exclude []string
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
		o.Exclude = ips
	}
}

// Option is IpRateLimit option.
type Option func(*options)

// RateLimiter  ip rate limiter
type RateLimiter struct {
	ips  map[string]*rate.Limiter
	mu   *sync.RWMutex
	opts options
}

func NewRateLimiter(opts ...Option) *RateLimiter {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}
	i := &RateLimiter{
		ips:  make(map[string]*rate.Limiter),
		mu:   &sync.RWMutex{},
		opts: opt,
	}

	return i
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *RateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.opts.limit, i.opts.burst)

	i.ips[ip] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise, calls AddIP to add IP address to the map
func (i *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	return limiter
}

func (i *RateLimiter) GetExcludes() []string {
	return i.opts.Exclude
}
