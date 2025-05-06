package auth

import (
	"bytes"
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/plus-core/core/v2/cache"
	"github.com/jxo-me/plus-core/core/v2/logger"
	"github.com/jxo-me/plus-core/pkg/v2/security"
	"io"
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
	key := r.Header.Get(v.config.APIKeyHeader)
	sig := r.Header.Get(v.config.SignatureHeader)
	tsStr := r.Header.Get(v.config.TimestampHeader)
	traceID := r.Header.Get(v.traceKey)
	clientIP := r.RemoteAddr
	ctx := r.GetCtx()

	if v.config.EnableIPBan && v.isBanned(ctx, clientIP) {
		v.log(ctx, "ip banned", traceID)
		return ErrIPBanned
	}

	if key == "" || sig == "" || tsStr == "" {
		if v.config.EnableIPBan {
			v.trackFailure(ctx, clientIP)
		}
		v.log(ctx, "missing headers", traceID)
		return ErrMissingHeader
	}

	ts, err := security.ParseTimestamp(tsStr)
	if err != nil || !security.IsTimestampValid(ts) {
		if v.config.EnableIPBan {
			v.trackFailure(ctx, clientIP)
		}
		v.log(ctx, "invalid timestamp", traceID)
		return ErrInvalidTimestamp
	}

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))

	secret, err := v.getSecretFromKey(ctx, key)
	if err != nil {
		if v.config.EnableIPBan {
			v.trackFailure(ctx, clientIP)
		}
		v.log(ctx, "invalid apiKey", traceID)
		return ErrInvalidAPIKey
	}

	signingString := security.BuildSigningString(r.Method, r.URL.Path, r.URL.Query(), body, ts)
	v.Debug(ctx, "signingString:", signingString)
	v.Debug(ctx, "secret:", secret)
	v.Debug(ctx, "sig:", sig)
	if v.config.EnableReplayCheck && v.isReplay(ctx, sig, ts) {
		if v.config.EnableIPBan {
			v.trackFailure(ctx, clientIP)
		}
		v.log(ctx, "replay attack", traceID)
		return ErrReplayAttack
	}

	if !v.strategy.Verify(signingString, secret, sig) {
		if v.config.EnableIPBan {
			v.trackFailure(ctx, clientIP)
		}
		v.log(ctx, "signature mismatch", traceID)
		return ErrSignatureMismatch
	}

	return nil
}

func (v *Verifier) getSecretFromKey(ctx context.Context, apiKey string) (string, error) {
	key := v.config.CachePrefix + "apiKey:" + apiKey
	// 优先缓存
	if v, err := v.cache.Get(ctx, key); err == nil {
		if v.String() != "" {
			return v.String(), nil
		}
	}
	// 自定义查库逻辑
	if v.config.GetSecretFunc != nil {
		secret, err := v.config.GetSecretFunc(apiKey)
		if err != nil {
			return "", err
		}
		if secret != "" {
			if err = v.cache.Set(ctx, key, secret, int(v.config.CacheTTL)); err != nil {
				v.logger.Errorf(ctx, "cache set SecretFromKey error for %s: %v", apiKey, err)
			}
			return secret, nil
		}
	}
	return "", ErrInvalidAPIKey
}

func (v *Verifier) Debug(ctx context.Context, key string, val ...interface{}) {
	if v.logger != nil {
		v.logger.Debugf(ctx, "auth debug [%s] val=%s", key, val)
	}
}

func (v *Verifier) log(ctx context.Context, reason, traceID string) {
	if v.logger != nil {
		v.logger.Debugf(ctx, "auth failed [%s] traceID=%s", reason, traceID)
	}
}

func (v *Verifier) isReplay(ctx context.Context, sig string, ts int64) bool {
	key := v.config.CachePrefix + "sig:" + sig

	exists, _ := v.cache.Get(ctx, key)
	if exists.Int() == 1 {
		return true
	}
	_ = v.cache.Set(ctx, key, 1, int(v.config.ReplayTTL))
	return false
}

func (v *Verifier) trackFailure(ctx context.Context, ip string) {
	key := v.config.CachePrefix + "ban:" + ip

	failures, _ := v.cache.Increase(ctx, key)
	if failures == 1 {
		_ = v.cache.Expire(ctx, key, v.config.BanExpire)
	}
	if failures >= 5 {
		_ = v.cache.Set(ctx, key+":block", 1, int(v.config.BanDuration))
	}
}

func (v *Verifier) isBanned(ctx context.Context, ip string) bool {
	val, _ := v.cache.Get(ctx, v.config.CachePrefix+"ban:"+ip+":block")
	return val.String() == "1"
}
