package auth

import (
	"github.com/jxo-me/plus-core/core/v2/cache"
	"github.com/jxo-me/plus-core/core/v2/logger"
)

type Config struct {
	Cache             cache.ICache
	Logger            logger.ILogger
	TraceHeaderKey    string
	CachePrefix       string
	SignatureStrategy SignatureStrategy
}

// SignatureStrategy 抽象签名算法策略接口（支持 HMAC、RSA 等）
type SignatureStrategy interface {
	Generate(signingString, secret string) string
	Verify(signingString, secret, givenSig string) bool
}
