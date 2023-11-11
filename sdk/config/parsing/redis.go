package parsing

import "github.com/jxo-me/plus-core/sdk/v2/config"

func ParseRedis(cfg *config.Redis) (chain.IChainer, error) {
	if cfg == nil {
		return nil, nil
	}

	return c, nil
}
