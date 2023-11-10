package config

import "github.com/jxo-me/plus-core/pkg/v2/security"

type Crypto struct {
	Enable    bool                      `json:"enable" yaml:"enable"`
	Algorithm string                    `json:"algorithm" yaml:"algorithm"`
	Rc4       security.Rc4CipherConfig  `json:"rc4" yaml:"rc4"`
	Rsa       security.RsaCiphersConfig `json:"rsa" yaml:"rsa"`
	Aes       security.AesCipherConfig  `json:"aes" yaml:"aes"`
}

type Metrics struct {
	Enable          bool      `json:"enable" yaml:"enable"`
	Path            string    `json:"path" yaml:"path"`
	SlowTime        int32     `json:"slowTime" yaml:"slowTime"`
	RequestDuration []float64 `json:"requestDuration" yaml:"requestDuration"`
}

type Jwt struct {
	Secret      string `yaml:"secret" json:"secret"`
	SigningKey  string `yaml:"signingKey" json:"signing_key"`
	Timeout     int64  `yaml:"timeout" json:"timeout"`
	MaxRefresh  int64  `yaml:"maxRefresh" json:"max_refresh"`
	IdentityKey string `yaml:"identityKey" json:"identity_key"`
}
