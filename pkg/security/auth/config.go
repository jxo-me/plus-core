package auth

import (
	"context"
	"github.com/jxo-me/plus-core/core/v2/cache"
	"github.com/jxo-me/plus-core/core/v2/logger"
	"github.com/jxo-me/plus-core/pkg/v2/security"
	"time"
)

type Config struct {
	Cache             cache.ICache                                                        // API Key -> Secret 缓存实现（支持 Redis）
	Logger            logger.ILogger                                                      // 日志记录接口（支持 traceID 日志输出）
	TraceHeaderKey    string                                                              // 请求中 traceID 所用的 Header 名（例如 Trace-Id）
	APIKeyHeader      string                                                              // 获取 API Key 的 Header 名
	SignatureHeader   string                                                              // 获取签名值的 Header 名
	TimestampHeader   string                                                              // 获取时间戳的 Header 名
	CachePrefix       string                                                              // Redis 缓存前缀
	SignatureStrategy security.SignatureStrategy                                          // 签名策略（如 HMAC 或 RSA）
	GetSecretFunc     func(ctx context.Context, apiKey string) (secret string, err error) // 自定义获取 API Secret 的方法
	EnableReplayCheck bool                                                                // 是否启用防重放签名校验
	EnableIPBan       bool                                                                // 是否启用 IP 封禁机制
	CacheTTL          time.Duration                                                       // API Key 缓存有效时间
	ReplayTTL         time.Duration                                                       // 重放签名有效期（用于标记是否已签过）
	BanDuration       time.Duration                                                       // 达到封禁阈值后封锁时长
	BanExpire         time.Duration                                                       // IP 连续失败计数的过期时间（计数器生命周期） (5 * time.Minute).Seconds()
}
