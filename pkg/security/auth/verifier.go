package auth

import (
	"bytes"
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/plus-core/core/v2/cache"
	"github.com/jxo-me/plus-core/core/v2/logger"
	"github.com/jxo-me/plus-core/pkg/v2/security"
	"io"
	"time"
)

type Verifier struct {
	config   *Config
	cache    cache.ICache
	logger   logger.ILogger
	traceKey string
	strategy security.SignatureStrategy
}

func NewVerifier(config *Config) *Verifier {
	return &Verifier{
		config:   config,
		cache:    config.Cache,
		logger:   config.Logger,
		traceKey: config.TraceHeaderKey,
		strategy: config.SignatureStrategy,
	}
}

func (v *Verifier) VerifyRequest(r *ghttp.Request) error {
	key := r.Header.Get("X-API-Key")
	sig := r.Header.Get("X-Signature")
	tsStr := r.Header.Get("Timestamp")
	traceID := r.Header.Get(v.traceKey)
	ctx := r.GetCtx()

	clientIP := r.RemoteAddr
	if v.isBanned(ctx, clientIP) {
		v.log(ctx, "ip banned", traceID)
		return ErrIPBanned
	}

	if key == "" || sig == "" || tsStr == "" {
		v.trackFailure(ctx, clientIP)
		v.log(ctx, "missing headers", traceID)
		return ErrMissingHeader
	}
	ts, err := security.ParseTimestamp(tsStr)
	if err != nil || !security.IsTimestampValid(ts) {
		v.trackFailure(ctx, clientIP)
		v.log(ctx, "invalid timestamp", traceID)
		return ErrInvalidTimestamp
	}

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))

	secret, err := v.getSecretFromKey(ctx, key)
	if err != nil {
		v.trackFailure(ctx, clientIP)
		v.log(ctx, "invalid apiKey", traceID)
		return ErrInvalidAPIKey
	}

	signingString := security.BuildSigningString(r.Method, r.URL.Path, r.URL.Query(), body, ts)

	if v.isReplay(ctx, sig, ts) {
		v.trackFailure(ctx, clientIP)
		v.log(ctx, "replay attack", traceID)
		return ErrReplayAttack
	}

	if !v.strategy.Verify(signingString, secret, sig) {
		v.trackFailure(ctx, clientIP)
		v.log(ctx, "signature mismatch", traceID)
		return ErrSignatureMismatch
	}

	return nil
}

func (v *Verifier) getSecretFromKey(ctx context.Context, apiKey string) (string, error) {
	if v.cache != nil {
		if v, err := v.cache.Get(ctx, apiKey); err == nil {
			return v.String(), nil
		}
	}
	return "", ErrInvalidAPIKey
}

func (v *Verifier) log(ctx context.Context, reason, traceID string) {
	if v.logger != nil {
		v.logger.Errorf(ctx, "auth failed [%s] traceID=%s", reason, traceID)
	}
}

func (v *Verifier) isReplay(ctx context.Context, sig string, ts int64) bool {
	key := v.config.CachePrefix + "sig:" + sig

	exists, _ := v.cache.Get(ctx, key)
	if exists.Int() == 1 {
		return true
	}
	_ = v.cache.Set(ctx, key, 1, int((5 * time.Minute).Seconds()))
	return false
}

func (v *Verifier) trackFailure(ctx context.Context, ip string) {
	key := v.config.CachePrefix + "ban:" + ip

	failures, _ := v.cache.Increase(ctx, key)
	if failures == 1 {
		_ = v.cache.Expire(ctx, key, 15*time.Minute)
	}
	if failures >= 5 {
		_ = v.cache.Set(ctx, key+":block", 1, int((30 * time.Minute).Seconds()))
	}
}

func (v *Verifier) isBanned(ctx context.Context, ip string) bool {
	val, _ := v.cache.Get(ctx, v.config.CachePrefix+"ban:"+ip+":block")
	return val.String() == "1"
}
